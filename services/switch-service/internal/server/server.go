package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/moorara/microservices-demo/services/switch-service/cmd/config"
	"github.com/moorara/microservices-demo/services/switch-service/cmd/version"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/internal/service"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	// GRPCServer is the interface for grpc.Server
	GRPCServer interface {
		Serve(net.Listener) error
		GracefulStop()
	}

	// Server is the server for services
	Server struct {
		port       string
		logger     *log.Logger
		grpcServer GRPCServer
		closers    []io.Closer
	}
)

// New creates a new server
func New(config config.Config) (*Server, error) {
	logger := log.NewLogger(config.ServiceName, "singleton", config.LogLevel)
	metrics := metrics.NewMetrics(config.ServiceName)

	sampler := trace.NewConstSampler()
	reporter := trace.NewReporter(config.JaegerLogSpans, config.JaegerAgentAddr)
	tracer, tracerCloser := trace.NewTracer(config.ServiceName, sampler, reporter, logger.Logger, metrics.Registry)

	switchService := service.NewSwitchService(logger, metrics, tracer)

	options := []grpc.ServerOption{}

	// Configure MTLS
	if config.CAChainFile != "" && config.ServerCertFile != "" && config.ServerKeyFile != "" {
		ca, err := ioutil.ReadFile(config.CAChainFile)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("Failed to append certificate authority")
		}

		cert, err := tls.LoadX509KeyPair(config.ServerCertFile, config.ServerKeyFile)
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
		port:       config.ServicePort,
		logger:     logger,
		grpcServer: grpcServer,
		closers:    []io.Closer{tracerCloser},
	}

	return server, nil
}

// Start starts listening
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	errs := make(chan error)

	// Handle OS signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs

		s.Shutdown()
		errs <- fmt.Errorf("%s", sig)
	}()

	// Listen for gRPC requests
	go func() {
		s.logger.Info(
			"version", version.Version,
			"revision", version.Revision,
			"branch", version.Branch,
			"buildTime", version.BuildTime,
			"message", fmt.Sprintf("gRPC server listening on port %s ...", s.port),
		)

		err := s.grpcServer.Serve(listener)
		if err != nil {
			s.logger.Error("message", fmt.Sprintf("gRPC server errored: %s", err.Error()))
			errs <- err
		}
	}()

	return <-errs
}

// Shutdown gracefully shutdowns the server
func (s *Server) Shutdown() error {
	// https://godoc.org/google.golang.org/grpc#Server.GracefulStop
	s.grpcServer.GracefulStop()

	for _, closer := range s.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	s.logger.Error(
		"message", "server was gracefully shutdown",
	)

	return nil
}
