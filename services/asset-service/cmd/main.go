package main

import (
	"fmt"

	"github.com/moorara/microservices-demo/services/asset-service/cmd/config"
	"github.com/moorara/microservices-demo/services/asset-service/cmd/server"
	"github.com/moorara/microservices-demo/services/asset-service/cmd/version"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/trace"
)

func main() {
	logger := log.NewLogger(config.Config.ServiceName, "singleton", config.Config.LogLevel)
	metrics := metrics.New(config.Config.ServiceName)

	// Tracer
	sampler := trace.NewConstSampler()
	reporter := trace.NewReporter(config.Config.JaegerLogSpans, config.Config.JaegerAgentAddr)
	tracer, tracerCloser := trace.NewTracer(config.Config.ServiceName, sampler, reporter, logger.Logger, metrics.Registry)
	defer tracerCloser.Close()

	server := server.New(config.Config.ServicePort, logger, metrics, tracer)

	logger.Info(
		"version", version.Version,
		"revision", version.Revision,
		"branch", version.Branch,
		"buildTime", version.BuildTime,
		"message", fmt.Sprintf("%s started.", config.Config.ServiceName),
	)

	server.Start()
}
