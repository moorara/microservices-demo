package trace

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func TestJaegerLogger(t *testing.T) {
	tests := []struct {
		name             string
		errorMsg         string
		infoMsg          string
		infoArgs         []interface{}
		expectedErrorMsg string
		expectedInfoMsg  string
	}{
		{
			name:             "",
			errorMsg:         "test error message",
			infoMsg:          "test %s %s",
			infoArgs:         []interface{}{"info", "message"},
			expectedErrorMsg: "test error message",
			expectedInfoMsg:  "test info message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Logger with pipe to read from
			rd, wr, _ := os.Pipe()
			dec := json.NewDecoder(rd)
			logger := log.NewJSONLogger(wr)

			jglogger := &jaegerLogger{logger}
			jglogger.Error(tc.errorMsg)
			jglogger.Infof(tc.infoMsg, tc.infoArgs...)

			var log map[string]interface{}

			// Verify Error
			dec.Decode(&log)
			assert.Equal(t, "error", log["level"])
			assert.Equal(t, tc.expectedErrorMsg, log["message"])

			// Verify Infof
			dec.Decode(&log)
			assert.Equal(t, "info", log["level"])
			assert.Equal(t, tc.expectedInfoMsg, log["message"])
		})
	}
}

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
			"WithLoggingAndMetrics",
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
