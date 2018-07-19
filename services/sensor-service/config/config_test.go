package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, defaultServiceName, Config.ServiceName)
	assert.Equal(t, defaultServicePort, Config.ServicePort)
	assert.Equal(t, defaultPostgresURL, Config.PostgresURL)
	assert.Equal(t, defaultLogLevel, Config.LogLevel)
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
		config := Spec{
			PostgresURL: tc.postgresURL,
		}
		assert.Equal(t, tc.expectedFullPostgresURL, config.GetFullPostgresURL())
	}
}
