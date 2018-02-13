package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/sensor-service/config"
	"github.com/moorara/microservices-demo/services/sensor-service/handler"
	"github.com/moorara/microservices-demo/services/sensor-service/middleware"
	"github.com/moorara/microservices-demo/services/sensor-service/service"
	"github.com/moorara/microservices-demo/services/sensor-service/util"
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
	metrics := util.NewMetrics("sensor_service")
	logger := util.NewLogger(config.LogLevel, config.ServiceName, "go-kit")

	metricsMiddleware := middleware.NewMetricsMiddleware(metrics)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)

	postgresDB := service.NewPostgresDB(config.PostgresURL)
	sensorHandler := handler.NewSensorHandler(postgresDB, logger)
	postSensorHandler := middleware.WrapAll(sensorHandler.PostSensor, metricsMiddleware, loggerMiddleware)
	getSensorsHandler := middleware.WrapAll(sensorHandler.GetSensors, metricsMiddleware, loggerMiddleware)
	getSensorHandler := middleware.WrapAll(sensorHandler.GetSensor, metricsMiddleware, loggerMiddleware)
	putSensorHandler := middleware.WrapAll(sensorHandler.PutSensor, metricsMiddleware, loggerMiddleware)
	deleteSensorHandler := middleware.WrapAll(sensorHandler.DeleteSensor, metricsMiddleware, loggerMiddleware)

	router := mux.NewRouter()
	router.NotFoundHandler = middleware.WrapAll(handler.GetNotFoundHandler(logger), loggerMiddleware)
	router.Methods("GET").Path("/health").HandlerFunc(handler.HealthHandler)
	router.Methods("GET").Path("/metrics").HandlerFunc(metrics.GetHandler().ServeHTTP)
	router.Methods("POST").Path("/v1/sensors").HandlerFunc(postSensorHandler)
	router.Methods("GET").Path("/v1/sensors").Queries("siteId", "{siteId}").HandlerFunc(getSensorsHandler)
	router.Methods("GET").Path("/v1/sensors/{id}").HandlerFunc(getSensorHandler)
	router.Methods("PUT").Path("/v1/sensors/{id}").HandlerFunc(putSensorHandler)
	router.Methods("DELETE").Path("/v1/sensors/{id}").HandlerFunc(deleteSensorHandler)

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
