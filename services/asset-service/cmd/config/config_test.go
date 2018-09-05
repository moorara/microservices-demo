package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name                    string
		EnvVars                 map[string]string
		expectedServiceName     string
		expectedServicePort     string
		expectedLogLevel        string
		expectedJaegerAgentAddr string
		expectedJaegerLogSpans  bool
	}{
		{
			name:                    "Defauts",
			EnvVars:                 map[string]string{},
			expectedServiceName:     defaultServiceName,
			expectedServicePort:     defaultServicePort,
			expectedLogLevel:        defaultLogLevel,
			expectedJaegerAgentAddr: defaultJaegerAgentAddr,
			expectedJaegerLogSpans:  defaultJaegerLogSpans,
		},
		{
			name: "Defauts",
			EnvVars: map[string]string{
				"SERVICE_NAME":      "test-service",
				"SERVICE_PORT":      ":5000",
				"LOG_LEVEL":         "debug",
				"JAEGER_AGENT_ADDR": "jaeger-agent:6831",
				"JAEGER_LOG_SPANS":  "true",
			},
			expectedServiceName:     "test-service",
			expectedServicePort:     ":5000",
			expectedLogLevel:        "debug",
			expectedJaegerAgentAddr: "jaeger-agent:6831",
			expectedJaegerLogSpans:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.EnvVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			assert.Equal(t, defaultServiceName, Config.ServiceName)
			assert.Equal(t, defaultServicePort, Config.ServicePort)
			assert.Equal(t, defaultLogLevel, Config.LogLevel)
			assert.Equal(t, defaultJaegerAgentAddr, Config.JaegerAgentAddr)
			assert.Equal(t, defaultJaegerLogSpans, Config.JaegerLogSpans)
		})
	}
}
