package main

import (
	"fmt"
	"math/rand"

	"github.com/moorara/microservices-demo/services/asset/cmd/config"
	"github.com/moorara/microservices-demo/services/asset/cmd/server"
	"github.com/moorara/microservices-demo/services/asset/cmd/version"
	"github.com/moorara/microservices-demo/services/asset/internal/db"
	"github.com/moorara/microservices-demo/services/asset/internal/queue"
	"github.com/moorara/microservices-demo/services/asset/internal/service"
	"github.com/moorara/microservices-demo/services/asset/internal/transport"
	"github.com/moorara/microservices-demo/services/asset/pkg/log"
	"github.com/moorara/microservices-demo/services/asset/pkg/metrics"
	"github.com/moorara/microservices-demo/services/asset/pkg/trace"
)

func main() {
	logger := log.NewLogger(config.Global.ServiceName, "singleton", config.Global.LogLevel)
	metrics := metrics.New(config.Global.ServiceName)

	// Tracer
	sampler := trace.NewConstSampler()
	reporter := trace.NewReporter(config.Global.JaegerLogSpans, config.Global.JaegerAgentAddr)
	tracer, tracerCloser := trace.NewTracer(config.Global.ServiceName, sampler, reporter, logger.Logger, metrics.Registry)
	defer tracerCloser.Close()

	// NATS Connection
	clientName := fmt.Sprintf("%s-%d", config.Global.ServiceName, rand.Int())
	conn, err := queue.NewNATSConnection(config.Global.NatsServers, clientName, config.Global.NatsUser, config.Global.NatsPassword)
	if err != nil {
		panic(err)
	}

	// CockroachDB ORM
	orm, err := db.NewCockroachORM(config.Global.CockroachAddr, config.Global.CockroachUser, config.Global.CockroachPassword, config.Global.CockroachDatabase, logger)
	if err != nil {
		panic(err)
	}

	alarmService := service.NewAlarmService(orm, logger, metrics, tracer)
	cameraService := service.NewCameraService(orm, logger, metrics, tracer)

	natsTransport := transport.NewNATSTransport(logger, metrics, tracer, conn, alarmService, cameraService)
	server := server.New(config.Global.ServicePort, natsTransport, logger, metrics)

	logger.Info(
		"version", version.Version,
		"revision", version.Revision,
		"branch", version.Branch,
		"goVersion", version.GoVersion,
		"buildTool", version.BuildTool,
		"buildTime", version.BuildTime,
		"message", fmt.Sprintf("%s started.", config.Global.ServiceName),
	)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
