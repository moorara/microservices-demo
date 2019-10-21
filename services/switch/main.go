package main

import (
	"fmt"

	"github.com/moorara/konfig"
	"github.com/moorara/microservices-demo/services/switch/cmd/config"
	"github.com/moorara/microservices-demo/services/switch/cmd/server"
	"github.com/moorara/microservices-demo/services/switch/cmd/version"
	"github.com/moorara/microservices-demo/services/switch/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch/pkg/log"
	"github.com/moorara/microservices-demo/services/switch/pkg/trace"
)

func main() {
	config := config.New()
	if err := konfig.Pick(&config); err != nil {
		panic(err)
	}

	logger := log.NewLogger(config.ServiceName, "singleton", config.LogLevel)
	metrics := metrics.New(config.ServiceName)

	sampler := trace.NewConstSampler()
	reporter := trace.NewReporter(config.JaegerLogSpans, config.JaegerAgentAddr)
	tracer, tracerCloser := trace.NewTracer(config.ServiceName, sampler, reporter, logger.Logger, metrics.Registry)
	defer tracerCloser.Close()

	server, err := server.New(config, logger, metrics, tracer)
	if err != nil {
		panic(err)
	}

	logger.Info(
		"version", version.Version,
		"revision", version.Revision,
		"branch", version.Branch,
		"goVersion", version.GoVersion,
		"buildTool", version.BuildTool,
		"buildTime", version.BuildTime,
		"message", fmt.Sprintf("%s started.", config.ServiceName),
	)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
