package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := New()

	assert.Equal(t, defaultServiceName, config.ServiceName)
	assert.Equal(t, defaultServicePort, config.ServicePort)
	assert.Equal(t, defaultLogLevel, config.LogLevel)
	assert.Equal(t, defaultPostgresHost, config.PostgresHost)
	assert.Equal(t, defaultPostgresPort, config.PostgresPort)
	assert.Equal(t, defaultPostgresDatabase, config.PostgresDatabase)
	assert.Equal(t, defaultPostgresUsername, config.PostgresUsername)
	assert.Equal(t, defaultPostgresPassword, config.PostgresPassword)
	assert.Equal(t, defaultJaegerAgentAddr, config.JaegerAgentAddr)
	assert.Equal(t, defaultJaegerLogSpans, config.JaegerLogSpans)
}
