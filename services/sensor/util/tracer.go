package util

import (
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/moorara/microservices-demo/services/sensor/config"
	"github.com/prometheus/client_golang/prometheus"

	opentracing "github.com/opentracing/opentracing-go"
	// jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"
	jaegerPrometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

type jaegerLogger struct {
	logger log.Logger
}

func (l *jaegerLogger) Error(msg string) {
	_ = level.Error(l.logger).Log("message", msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	_ = l.logger.Log("message", fmt.Sprintf(msg, args...))
}

// NewTracer creates a new tracer
func NewTracer(config config.Config, logger log.Logger, registerer prometheus.Registerer) (opentracing.Tracer, io.Closer) {
	jgConfig := &jaegerConfig.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: config.JaegerAgentAddr,
			LogSpans:           config.JaegerLogSpans,
		},
	}

	// jglogger := jaeger.NullLogger
	jglogger := &jaegerLogger{logger}
	loggerOption := jaegerConfig.Logger(jglogger)

	// metricsFactory := jaegerMetrics.NullFactory
	registererOption := jaegerPrometheus.WithRegisterer(registerer)
	metricsFactory := jaegerPrometheus.New(registererOption).Namespace(jaegerMetrics.NSOptions{Name: config.ServiceName})
	metricsOption := jaegerConfig.Metrics(metricsFactory)

	jaegerTracer, closer, err := jgConfig.NewTracer(loggerOption, metricsOption)
	if err != nil {
		panic(err)
	}

	return jaegerTracer, closer
}
