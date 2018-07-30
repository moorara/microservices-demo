package main

import (
	"context"
	"fmt"
	"time"

	gotoConfig "github.com/moorara/goto/config"
	"github.com/moorara/microservices-demo/services/switch-service/cmd/config"
	"github.com/moorara/microservices-demo/services/switch-service/cmd/version"
	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/server"
	"github.com/moorara/microservices-demo/services/switch-service/internal/service"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/trace"
)

const (
	timeout = 30 * time.Second
)

func main() {
	config := config.New()
	gotoConfig.Pick(&config)

	logger := log.NewLogger(config.ServiceName, "singleton", config.LogLevel)
	metrics := metrics.New(config.ServiceName)

	sampler := trace.NewConstSampler()
	reporter := trace.NewReporter(config.JaegerLogSpans, config.JaegerAgentAddr)
	tracer, tracerCloser := trace.NewTracer(config.ServiceName, sampler, reporter, logger.Logger, metrics.Registry)
	defer tracerCloser.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	arangoService := service.NewArangoService()
	err := arangoService.Connect(
		ctx,
		config.ArangoEndpoints, config.ArangoUser, config.ArangoPassword,
		config.ArangoDatabase, config.ArangoCollection,
	)

	if err != nil {
		panic(err)
	}

	switchService := service.NewSwitchService(arangoService, logger, metrics, tracer)

	server, err := server.New(
		config.ServiceHTTPPort, config.ServiceGRPCPort,
		config.CAChainFile, config.ServerCertFile, config.ServerKeyFile,
		switchService, logger, metrics,
	)

	if err != nil {
		panic(err)
	}

	logger.Info(
		"version", version.Version,
		"revision", version.Revision,
		"branch", version.Branch,
		"buildTime", version.BuildTime,
		"message", fmt.Sprintf("%s started.", config.ServiceName),
	)

	server.Start()
}
