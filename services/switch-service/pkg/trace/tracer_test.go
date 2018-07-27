package trace

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func TestNewConstSampler(t *testing.T) {
	tests := []struct {
		name          string
		expectedParam float64
	}{
		{"Always", 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sampler := NewConstSampler()
			assert.Equal(t, sampler.Type, "const")
			assert.Equal(t, sampler.Param, tc.expectedParam)
		})
	}
}

func TestNewProbabilisticSampler(t *testing.T) {
	tests := []struct {
		name        string
		probability float64
	}{
		{"Never", 0.0},
		{"HalfTheTime", 0.5},
		{"Always", 1.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sampler := NewProbabilisticSampler(tc.probability)
			assert.Equal(t, sampler.Type, "probabilistic")
			assert.Equal(t, sampler.Param, tc.probability)
		})
	}
}

func TestNewRateLimitingSampler(t *testing.T) {
	tests := []struct {
		name           string
		spansPerSecond float64
	}{
		{"1SpansPerSecond", 1},
		{"10SpansPerSecond", 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sampler := NewRateLimitingSampler(tc.spansPerSecond)
			assert.Equal(t, sampler.Type, "rateLimiting")
			assert.Equal(t, sampler.Param, tc.spansPerSecond)
		})
	}
}

func TestNewReporter(t *testing.T) {
	tests := []struct {
		name            string
		logSpans        bool
		jaegerAgentAddr string
	}{
		{"UsingLocalAgent", false, ""},
		{"UsingNonLocalAgent", true, "jaeger-agent:6831"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reporter := NewReporter(tc.logSpans, tc.jaegerAgentAddr)
			assert.Equal(t, tc.logSpans, reporter.LogSpans)
			assert.Equal(t, tc.jaegerAgentAddr, reporter.LocalAgentHostPort)
		})
	}
}

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		sampler     *jaegerConfig.SamplerConfig
		reporter    *jaegerConfig.ReporterConfig
		logger      log.Logger
		registerer  prometheus.Registerer
	}{
		{
			"Default",
			"service_name",
			&jaegerConfig.SamplerConfig{},
			&jaegerConfig.ReporterConfig{},
			nil,
			nil,
		},
		{
			"WithLogging",
			"service_name",
			&jaegerConfig.SamplerConfig{},
			&jaegerConfig.ReporterConfig{},
			log.NewNopLogger(),
			nil,
		},
		{
			"WithMetrics",
			"service_name",
			&jaegerConfig.SamplerConfig{},
			&jaegerConfig.ReporterConfig{},
			nil,
			prometheus.NewRegistry(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer, closer := NewTracer(tc.serviceName, tc.sampler, tc.reporter, tc.logger, tc.registerer)
			defer closer.Close()

			assert.NotNil(t, tracer)
			assert.NotNil(t, closer)
		})
	}
}
