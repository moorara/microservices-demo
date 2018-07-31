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

	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/internal/transport"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
)

const (
	stopTimeout = 30 * time.Second
)

type (
	// Server coordinates http and grpc transports
	Server interface {
		Start() error
		Stop()
	}

	server struct {
		grpcPort   string
		logger     *log.Logger
		httpServer transport.HTTPServer
		grpcServer transport.GRPCServer
	}
)

// New creates a new server
func New(httpPort, grpcPort, caFile, certFile, keyFile string,
	switchService proto.SwitchServiceServer, logger *log.Logger, metrics *metrics.Metrics) (Server, error) {
	var err error
	s := &server{
		grpcPort: grpcPort,
		logger:   logger,
	}

	s.httpServer = transport.NewHTTPServer(httpPort, s.liveHandler, s.readyHandler, metrics.Handler().ServeHTTP)
	s.grpcServer, err = transport.NewGRPCServer(caFile, certFile, keyFile, switchService)
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
	w.WriteHeader(http.StatusOK)
}

func (s *server) Start() error {
	errs := make(chan error)
	done := make(chan struct{}, 1)

	listener, err := net.Listen("tcp", s.grpcPort)
	if err != nil {
		return err
	}

	// Handle OS signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-done:
		case sig := <-sigs:
			errs <- errors.New(sig.String())
		}
	}()

	// Listen for http requests
	go func() {
		s.logger.Info("message", "http server listening ...")

		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("message", fmt.Sprintf("http server errored: %s", err.Error()))
			errs <- err
		}
	}()

	// Listen for gRPC requests
	go func() {
		s.logger.Info("message", "gRPC server listening ...")

		err := s.grpcServer.Serve(listener)
		if err != nil {
			s.logger.Error("message", fmt.Sprintf("gRPC server errored: %s", err.Error()))
			errs <- err
		}
	}()

	err = <-errs
	done <- struct{}{}
	s.Stop()

	return err
}

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	defer cancel()

	s.httpServer.Shutdown(ctx)
	s.grpcServer.GracefulStop()

	s.logger.Info("message", "server was gracefully shutdown.")
}
