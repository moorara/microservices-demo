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
	assert.Equal(t, defaultJaegerAgentAddr, config.JaegerAgentAddr)
	assert.Equal(t, defaultJaegerLogSpans, config.JaegerLogSpans)
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
