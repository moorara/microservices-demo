package main

import (
	"fmt"
	"math/rand"

	. "github.com/moorara/microservices-demo/services/asset-service/cmd/config"
	"github.com/moorara/microservices-demo/services/asset-service/cmd/server"
	"github.com/moorara/microservices-demo/services/asset-service/cmd/version"
	"github.com/moorara/microservices-demo/services/asset-service/internal/db"
	"github.com/moorara/microservices-demo/services/asset-service/internal/queue"
	"github.com/moorara/microservices-demo/services/asset-service/internal/service"
	"github.com/moorara/microservices-demo/services/asset-service/internal/transport"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/trace"
)

func main() {
	logger := log.NewLogger(Config.ServiceName, "singleton", Config.LogLevel)
	metrics := metrics.New(Config.ServiceName)

	// Tracer
	sampler := trace.NewConstSampler()
	reporter := trace.NewReporter(Config.JaegerLogSpans, Config.JaegerAgentAddr)
	tracer, tracerCloser := trace.NewTracer(Config.ServiceName, sampler, reporter, logger.Logger, metrics.Registry)
	defer tracerCloser.Close()

	// NATS Connection
	clientName := fmt.Sprintf("%s-%d", Config.ServiceName, rand.Int())
	conn, err := queue.NewNATSConnection(Config.NatsServers, clientName, Config.NatsUser, Config.NatsPassword)
	if err != nil {
		panic(err)
	}

	// CockroachDB ORM
	orm, err := db.NewCockroachORM(Config.CockroachAddr, Config.CockroachUser, Config.CockroachPassword, Config.CockroachDatabase, logger)
	if err != nil {
		panic(err)
	}

	alarmService := service.NewAlarmService(orm, logger, metrics, tracer)
	cameraService := service.NewCameraService(orm, logger, metrics, tracer)

	natsTransport := transport.NewNATSTransport(logger, metrics, tracer, conn, alarmService, cameraService)
	server := server.New(Config.ServicePort, natsTransport, logger, metrics)

	logger.Info(
		"version", version.Version,
		"revision", version.Revision,
		"branch", version.Branch,
		"goVersion", version.GoVersion,
		"buildTool", version.BuildTool,
		"buildTime", version.BuildTime,
		"message", fmt.Sprintf("%s started.", Config.ServiceName),
	)

	server.Start()
}
