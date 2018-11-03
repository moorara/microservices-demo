package server

import (
	"context"
	"io"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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

	// HTTPServer represents a http server
	HTTPServer struct {
		config  config.Config
		logger  log.Logger
		closers []io.Closer
		server  Server
	}
)

// New creates a new http server
func New(config config.Config) *HTTPServer {
	metrics := util.NewMetrics("sensor_service")
	logger := util.NewLogger(config.LogLevel, config.ServiceName, "singleton")
	tracer, tracerCloser := util.NewTracer(config, logger, metrics.Registry)

	metricsMiddleware := middleware.NewMetricsMiddleware(metrics)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)
	tracerMiddleware := middleware.NewTracerMiddleware(tracer)

	postgresDB := service.NewPostgresDB(logger, config.PostgresHost, config.PostgresPort, config.PostgresDatabase, config.PostgresUsername, config.PostgresPassword)
	sensorHandler := handler.NewSensorHandler(postgresDB, logger, tracer)
	postSensorHandler := middleware.WrapAll(sensorHandler.PostSensor, metricsMiddleware, loggerMiddleware, tracerMiddleware)
	getSensorsHandler := middleware.WrapAll(sensorHandler.GetSensors, metricsMiddleware, loggerMiddleware, tracerMiddleware)
	getSensorHandler := middleware.WrapAll(sensorHandler.GetSensor, metricsMiddleware, loggerMiddleware, tracerMiddleware)
	putSensorHandler := middleware.WrapAll(sensorHandler.PutSensor, metricsMiddleware, loggerMiddleware, tracerMiddleware)
	deleteSensorHandler := middleware.WrapAll(sensorHandler.DeleteSensor, metricsMiddleware, loggerMiddleware, tracerMiddleware)

	router := mux.NewRouter()
	router.NotFoundHandler = middleware.WrapAll(handler.GetNotFoundHandler(logger), loggerMiddleware, tracerMiddleware)
	router.Methods("GET").Path("/health").HandlerFunc(handler.HealthHandler)
	router.Methods("GET").Path("/metrics").HandlerFunc(metrics.GetHandler().ServeHTTP)
	router.Methods("POST").Path("/v1/sensors").HandlerFunc(postSensorHandler)
	router.Methods("GET").Path("/v1/sensors").Queries("siteId", "{siteId}").HandlerFunc(getSensorsHandler)
	router.Methods("GET").Path("/v1/sensors/{id}").HandlerFunc(getSensorHandler)
	router.Methods("PUT").Path("/v1/sensors/{id}").HandlerFunc(putSensorHandler)
	router.Methods("DELETE").Path("/v1/sensors/{id}").HandlerFunc(deleteSensorHandler)

	return &HTTPServer{
		config:  config,
		logger:  logger,
		closers: []io.Closer{tracerCloser},
		server: &http.Server{
			Addr:    config.ServicePort,
			Handler: router,
		},
	}
}

// Start starts the server
func (s *HTTPServer) Start() error {
	s.logger.Log("message", "Listening on port "+s.config.ServicePort+" ...")
	level.Debug(s.logger).Log(
		"message", "configuration values",
		"config.postgresHost", s.config.PostgresHost,
		"config.postgresPort", s.config.PostgresPort,
		"config.postgresDatabase", s.config.PostgresDatabase,
		"config.postgresUsername", s.config.PostgresUsername,
		"config.postgresPassword", s.config.PostgresPassword,
		"config.jaegerAgentAddr", s.config.JaegerAgentAddr,
		"config.jaegerLogSpans", s.config.JaegerLogSpans,
	)

	return s.server.ListenAndServe()
}

// Close implements io.Closer to free up all resources used by server
func (s *HTTPServer) Close() error {
	for _, closer := range s.closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
