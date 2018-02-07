package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/go-service/config"
	"github.com/moorara/microservices-demo/services/go-service/handler"
	"github.com/moorara/microservices-demo/services/go-service/middleware"
	"github.com/moorara/microservices-demo/services/go-service/service"
	"github.com/moorara/microservices-demo/services/go-service/util"
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

	postgresDB := service.NewPostgresDB(config.PostgresURL)
	voteHandler := handler.NewVoteHandler(postgresDB, logger)
	postVoteHandler := middleware.WrapAll(voteHandler.PostVote, metricsMiddleware, loggerMiddleware)
	getVotesHandler := middleware.WrapAll(voteHandler.GetVotes, metricsMiddleware, loggerMiddleware)
	getVoteHandler := middleware.WrapAll(voteHandler.GetVote, metricsMiddleware, loggerMiddleware)
	deleteVoteHandler := middleware.WrapAll(voteHandler.DeleteVote, metricsMiddleware, loggerMiddleware)

	router := mux.NewRouter()
	router.NotFoundHandler = middleware.WrapAll(handler.GetNotFoundHandler(logger), loggerMiddleware)
	router.Methods("GET").Path("/health").HandlerFunc(handler.HealthHandler)
	router.Methods("GET").Path("/metrics").HandlerFunc(metrics.GetHandler().ServeHTTP)
	router.Methods("POST").Path("/v1/votes").HandlerFunc(postVoteHandler)
	router.Methods("GET").Path("/v1/votes").Queries("linkId", "{linkId}").HandlerFunc(getVotesHandler)
	router.Methods("GET").Path("/v1/votes/{id}").HandlerFunc(getVoteHandler)
	router.Methods("DELETE").Path("/v1/votes/{id}").HandlerFunc(deleteVoteHandler)

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
