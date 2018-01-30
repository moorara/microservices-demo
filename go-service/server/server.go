package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/toys/microservices/go-service/config"
	"github.com/moorara/toys/microservices/go-service/handler"
	"github.com/moorara/toys/microservices/go-service/middleware"
	"github.com/moorara/toys/microservices/go-service/service"
	"github.com/moorara/toys/microservices/go-service/util"
)

type (
	// Server represents a generic server
	Server interface {
		ListenAndServe() error
		Shutdown(context.Context) error
	}

	// HTTPServer representa a http server
	HTTPServer struct {
		config config.Config
		logger log.Logger
		server Server
	}
)

// New creates a new http server
func New(config config.Config) *HTTPServer {
	metrics := util.NewMetrics("go_service")
	logger := util.NewLogger(config.LogLevel, config.ServiceName, "go-kit")

	metricsMiddleware := middleware.NewMetricsMiddleware(metrics)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)

	redisPersister := service.NewRedisPersister(config.RedisURL)
	voteHandler := handler.NewVoteHandler(redisPersister, logger)
	postVoteHandler := middleware.WrapAll(voteHandler.PostVote, metricsMiddleware, loggerMiddleware)
	getVoteHandler := middleware.WrapAll(voteHandler.GetVote, metricsMiddleware, loggerMiddleware)

	router := mux.NewRouter()
	router.NotFoundHandler = handler.GetNotFoundHandler(logger)
	router.HandleFunc("/health", handler.HealthHandler).Methods("GET")
	router.HandleFunc("/metrics", metrics.GetHandler().ServeHTTP)
	router.HandleFunc("/v1/votes", postVoteHandler).Methods("POST")
	router.HandleFunc("/v1/votes/{id}", getVoteHandler).Methods("GET")

	return &HTTPServer{
		config: config,
		logger: logger,
		server: &http.Server{
			Addr:    config.ServicePort,
			Handler: router,
		},
	}
}

// Start starts the server
func (s *HTTPServer) Start() error {
	s.logger.Log("message", "Listening on port "+s.config.ServicePort+" ...")
	return s.server.ListenAndServe()
}
