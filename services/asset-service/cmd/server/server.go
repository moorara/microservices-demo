package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/asset-service/internal/transport"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/opentracing/opentracing-go"
)

type (
	// HTTPServer is the interface for http.Server
	HTTPServer interface {
		ListenAndServe() error
		Shutdown(context.Context) error
	}

	// Server manages a http.Server
	Server struct {
		logger        *log.Logger
		httpServer    HTTPServer
		natsTransport transport.NATSTransport
	}
)

// New creates a new Server
func New(port string, natsTransport transport.NATSTransport, logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) *Server {
	router := mux.NewRouter()
	server := &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:    port,
			Handler: router,
		},
		natsTransport: natsTransport,
	}

	router.NotFoundHandler = http.HandlerFunc(server.notFound)
	router.Methods("GET").Path("/liveness").HandlerFunc(server.liveness)
	router.Methods("GET").Path("/readiness").HandlerFunc(server.readiness)
	router.Methods("GET").Path("/metrics").Handler(metrics.Handler())

	// monitorMiddleware := middleware.NewMonitorMiddleware(logger, metrics, tracer)
	// graphql = monitorMiddleware.Wrap(graphql)
	// router.Path("/graphql").HandlerFunc(graphql)

	return server
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.logger.Warn("message", "Not found.", "method", r.Method, "url", r.URL.Path)
}

func (s *Server) liveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) readiness(w http.ResponseWriter, _ *http.Request) {
	// TODO should respond with 200 only if the service is ready to serve requests!
	w.WriteHeader(http.StatusOK)
}

// Start starts the http server!
func (s *Server) Start() error {
	errs := make(chan error)
	done := make(chan struct{}, 1)

	// Capture OS signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-done:
		case sig := <-sigs:
			errs <- fmt.Errorf("Interrupted by signal %s", sig.String())
		}
	}()

	// Listen for http requests
	go func() {
		s.logger.Info("message", "http server listening ...")
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("message", "http server errored.", "error", err)
			errs <- err
		}
	}()

	err := <-errs
	done <- struct{}{}
	s.Stop()

	return err
}

// Stop stops the http server!
func (s *Server) Stop() {
	timeout := time.Duration(30 * time.Second) // Default Kubernetes grace period
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.httpServer.Shutdown(ctx)
	s.logger.Info("message", "server was gracefully shutdown.")
}
