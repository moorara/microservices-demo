package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	// HTTPServer is the interface for http.Server
	HTTPServer interface {
		ListenAndServe() error
		Shutdown(context.Context) error
	}

	// GRPCServer is the interface for grpc.Server
	GRPCServer interface {
		Serve(net.Listener) error
		GracefulStop()
	}

	// Server is the server for services
	Server struct {
		logger     *log.Logger
		httpPort   string
		httpServer HTTPServer
		grpcPort   string
		grpcServer GRPCServer
		ready      bool
	}
)

// New creates a new server
func New(
	httpPort, grpcPort string,
	caChainFile, serverCertFile, serverKeyFile string,
	switchService proto.SwitchServiceServer, logger *log.Logger, metrics *metrics.Metrics,
) (*Server, error) {

	options := []grpc.ServerOption{}

	// Configure MTLS
	if caChainFile != "" && serverCertFile != "" && serverKeyFile != "" {
		ca, err := ioutil.ReadFile(caChainFile)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("Failed to append certificate authority")
		}

		cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    pool,
		}

		creds := credentials.NewTLS(tlsConfig)
		options = append(options, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(options...)
	proto.RegisterSwitchServiceServer(grpcServer, switchService)

	server := &Server{
		logger:     logger,
		httpPort:   httpPort,
		httpServer: nil,
		grpcPort:   grpcPort,
		grpcServer: grpcServer,
	}

	httpRouter := mux.NewRouter()
	httpRouter.NotFoundHandler = http.NotFoundHandler()
	httpRouter.Methods("GET").Path("/live").HandlerFunc(server.LiveHandler)
	httpRouter.Methods("GET").Path("/ready").HandlerFunc(server.ReadyHandler)
	httpRouter.Methods("GET").Path("/metrics").Handler(metrics.Handler())

	server.httpServer = &http.Server{
		Addr:    httpPort,
		Handler: httpRouter,
	}

	return server, nil
}

// LiveHandler implements liveness prob
func (s *Server) LiveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ReadyHandler implements readiness prob
func (s *Server) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if s.ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// Start starts listening
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.grpcPort)
	if err != nil {
		return err
	}

	errs := make(chan error)

	// Handle OS signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs

		errs <- errors.New(sig.String())
	}()

	// Listen for http requests
	go func() {
		s.logger.Info("message", fmt.Sprintf("http server listening on port %s ...", s.httpPort))

		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("message", fmt.Sprintf("http server errored: %s", err.Error()))
			errs <- err
		}
	}()

	// Listen for gRPC requests
	go func() {
		s.logger.Info("message", fmt.Sprintf("gRPC server listening on port %s ...", s.grpcPort))

		// Determine if service is ready!
		s.ready = true
		err := s.grpcServer.Serve(listener)
		s.ready = false

		if err != nil {
			s.logger.Error("message", fmt.Sprintf("gRPC server errored: %s", err.Error()))
			errs <- err
		}
	}()

	err = <-errs
	s.Shutdown()

	return err
}

// Shutdown gracefully shutdowns the server
func (s *Server) Shutdown() {
	ctx := context.Background()

	s.httpServer.Shutdown(ctx)  // https://godoc.org/net/http#Server.Shutdown
	s.grpcServer.GracefulStop() // https://godoc.org/google.golang.org/grpc#Server.GracefulStop

	s.logger.Info("message", "server was gracefully shutdown.")
}
