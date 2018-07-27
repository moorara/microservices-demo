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
	assert.Equal(t, defaultJaegerAgentAddr, config.JaegerAgentAddr)
	assert.Equal(t, defaultJaegerLogSpans, config.JaegerLogSpans)
	assert.Empty(t, config.CAChainFile)
	assert.Empty(t, config.ServerCertFile)
	assert.Empty(t, config.ServerKeyFile)
}
