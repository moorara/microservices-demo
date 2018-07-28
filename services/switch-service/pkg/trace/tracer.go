package trace

import (
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"
	jaegerPrometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

// jaegerLogger implements jaeger.Logger
type jaegerLogger struct {
	logger log.Logger
}

func (l *jaegerLogger) Error(msg string) {
	level.Error(l.logger).Log("message", msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	level.Info(l.logger).Log("message", fmt.Sprintf(msg, args...))
}

// NewConstSampler creates a constant Jaeger sampler
func NewConstSampler() *jaegerConfig.SamplerConfig {
	return &jaegerConfig.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
}

// NewProbabilisticSampler creates a probabilistic Jaeger sampler
func NewProbabilisticSampler(probability float64) *jaegerConfig.SamplerConfig {
	return &jaegerConfig.SamplerConfig{
		Type:  "probabilistic",
		Param: probability,
	}
}

// NewRateLimitingSampler creates a rate limiting Jaeger sampler
func NewRateLimitingSampler(spansPerSecond float64) *jaegerConfig.SamplerConfig {
	return &jaegerConfig.SamplerConfig{
		Type:  "rateLimiting",
		Param: spansPerSecond,
	}
}

// NewReporter creates a Jaeger reporter
func NewReporter(logSpans bool, jaegerAgentAddr string) *jaegerConfig.ReporterConfig {
	return &jaegerConfig.ReporterConfig{
		LogSpans:           logSpans,
		LocalAgentHostPort: jaegerAgentAddr,
	}
}

// NewTracer creates a new tracer
func NewTracer(serviceName string, sampler *jaegerConfig.SamplerConfig, reporter *jaegerConfig.ReporterConfig, logger log.Logger, registerer prometheus.Registerer) (opentracing.Tracer, io.Closer) {
	jgConfig := &jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler:     sampler,
		Reporter:    reporter,
	}

	var jglogger jaeger.Logger
	if logger == nil {
		jglogger = jaeger.NullLogger
	} else {
		jglogger = &jaegerLogger{logger}
	}
	loggerOption := jaegerConfig.Logger(jglogger)

	var metricsFactory jaegerMetrics.Factory
	if registerer == nil {
		metricsFactory = jaegerMetrics.NullFactory
	} else {
		registererOption := jaegerPrometheus.WithRegisterer(registerer)
		metricsFactory = jaegerPrometheus.New(registererOption)
		metricsFactory = metricsFactory.Namespace(serviceName, map[string]string{})
	}
	metricsOption := jaegerConfig.Metrics(metricsFactory)

	jaegerTracer, closer, err := jgConfig.NewTracer(loggerOption, metricsOption)
	if err != nil {
		panic(err)
	}

	return jaegerTracer, closer
}
