package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := New()

	assert.Equal(t, defaultServiceName, config.ServiceName)
	assert.Equal(t, defaultServicePort, config.ServicePort)
	assert.Equal(t, defaultPostgresURL, config.PostgresURL)
	assert.Equal(t, defaultLogLevel, config.LogLevel)
	assert.Equal(t, defaultJaegerAgentHost, config.JaegerAgentHost)
	assert.Equal(t, defaultJaegerAgentPort, config.JaegerAgentPort)
	assert.Equal(t, defaultJaegerReporterLogSpans, config.JaegerReporterLogSpans)
}

func TestGetFullPostgresURL(t *testing.T) {
	tests := []struct {
		postgresURL             string
		expectedFullPostgresURL string
	}{
		{"postgres://root@localhost", "postgres://root@localhost" + "/" + dbName + dbOpts},
		{"postgres://root@postgres", "postgres://root@postgres" + "/" + dbName + dbOpts},
	}

	for _, tc := range tests {
		config := Config{
			PostgresURL: tc.postgresURL,
		}
		assert.Equal(t, tc.expectedFullPostgresURL, config.GetFullPostgresURL())
	}
}

func TestGetJaegerAgentURL(t *testing.T) {
	tests := []struct {
		jaegerAgentHost        string
		jaegerAgentPort        int
		expectedJaegerAgentURL string
	}{
		{"localhost", 6831, "localhost:6831"},
		{"jaeger-agent", 6831, "jaeger-agent:6831"},
	}

	for _, tc := range tests {
		config := Config{
			JaegerAgentHost: tc.jaegerAgentHost,
			JaegerAgentPort: tc.jaegerAgentPort,
		}
		assert.Equal(t, tc.expectedJaegerAgentURL, config.GetJaegerAgentURL())
	}
}
