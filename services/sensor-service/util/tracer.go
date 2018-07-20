package util

import (
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/moorara/microservices-demo/services/sensor-service/config"

	opentracing "github.com/opentracing/opentracing-go"
	// jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"
	// jaegerPrometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

type jaegerLogger struct {
	logger log.Logger
}

func (l *jaegerLogger) Error(msg string) {
	level.Error(l.logger).Log("message", msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	l.logger.Log("message", fmt.Sprintf(msg, args...))
}

// NewTracer creates a new tracer
func NewTracer(config config.Config, logger log.Logger) (opentracing.Tracer, io.Closer) {
	jgConfig := &jaegerConfig.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: config.GetJaegerAgentURL(),
			LogSpans:           config.JaegerReporterLogSpans,
		},
	}

	// jglogger := jaeger.NullLogger
	jglogger := &jaegerLogger{logger}
	loggerOption := jaegerConfig.Logger(jglogger)

	metricsFactory := jaegerMetrics.NullFactory
	// regOption := jaegerPrometheus.WithRegisterer(registerer)
	// metricsFactory := jaegerPrometheus.New(regOption)
	metricsOption := jaegerConfig.Metrics(metricsFactory)

	jaegerTracer, closer, err := jgConfig.NewTracer(loggerOption, metricsOption)
	if err != nil {
		panic(err)
	}

	return jaegerTracer, closer
}
