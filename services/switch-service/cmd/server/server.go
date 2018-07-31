package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moorara/microservices-demo/services/switch-service/cmd/config"
	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/service"
	"github.com/moorara/microservices-demo/services/switch-service/internal/transport"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/opentracing/opentracing-go"
)

type (
	// Server coordinates http and grpc transports
	Server interface {
		Start() error
		Stop()
	}

	server struct {
		ready         bool
		config        config.Config
		logger        *log.Logger
		arangoService service.ArangoService
		httpServer    transport.HTTPServer
		grpcServer    transport.GRPCServer
	}
)

// New creates a new server
func New(config config.Config, logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) (Server, error) {
	var err error
	s := &server{
		config: config,
		logger: logger,
	}

	s.arangoService, err = service.NewArangoService(config.ArangoEndpoints, config.ArangoUser, config.ArangoPassword)
	if err != nil {
		return nil, err
	}

	s.httpServer = transport.NewHTTPServer(config.ServiceHTTPPort, s.liveHandler, s.readyHandler, metrics.Handler().ServeHTTP)

	switchService := service.NewSwitchService(s.arangoService, logger, metrics, tracer)
	s.grpcServer, err = transport.NewGRPCServer(config.CAChainFile, config.ServerCertFile, config.ServerKeyFile, switchService)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// liveHandler implements liveness prob
func (s *server) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readyHandler implements readiness prob
func (s *server) readyHandler(w http.ResponseWriter, r *http.Request) {
	if s.ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *server) Start() error {
	errs := make(chan error)
	sigDone := make(chan struct{}, 1)
	connDone := make(chan struct{}, 1)

	timeout := time.Duration(s.config.ServerTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Handle OS signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-sigDone:
		case sig := <-sigs:
			errs <- errors.New(sig.String())
		}
	}()

	// Connect to database
	go func() {
		for {
			err := s.arangoService.Connect(ctx, s.config.ArangoDatabase, s.config.ArangoCollection)
			if err == nil {
				s.logger.Info("message", "Connected to database.")
				s.ready = true
				return
			}

			select {
			case <-time.After(time.Second):
			case <-connDone:
				return
			case <-ctx.Done():
				s.logger.Error("message", "Failed to connect to database.")
				errs <- err
				return
			}
		}
	}()

	// Listen for http requests
	go func() {
		s.logger.Info("message", "http server listening ...")

		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("message", fmt.Sprintf("http server errored: %s", err))
			errs <- err
		}
	}()

	// Listen for gRPC requests
	go func() {
		listener, err := net.Listen("tcp", s.config.ServiceGRPCPort)
		if err != nil {
			s.logger.Error("message", "Failed to listen on grpc port.")
			errs <- err
			return
		}

		s.logger.Info("message", "gRPC server listening ...")

		err = s.grpcServer.Serve(listener)
		if err != nil {
			s.logger.Error("message", fmt.Sprintf("gRPC server errored: %s", err))
			errs <- err
		}
	}()

	err := <-errs
	sigDone <- struct{}{}
	connDone <- struct{}{}
	s.Stop()

	return err
}

func (s *server) Stop() {
	timeout := time.Duration(s.config.ServerTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.httpServer.Shutdown(ctx)
	s.grpcServer.GracefulStop()

	s.logger.Info("message", "server was gracefully shutdown.")
}
