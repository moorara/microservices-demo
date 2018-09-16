package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name                      string
		expectedLogLevel          string
		expectedServiceName       string
		expectedServicePort       string
		expectedNatsServers       []string
		expectedNatsUser          string
		expectedNatsPassword      string
		expectedCockroachAddr     string
		expectedCockroachUser     string
		expectedCockroachPassword string
		expectedCockroachDatabase string
		expectedJaegerAgentAddr   string
		expectedJaegerLogSpans    bool
	}{
		{
			name:                      "Defauts",
			expectedLogLevel:          defaultLogLevel,
			expectedServiceName:       defaultServiceName,
			expectedServicePort:       defaultServicePort,
			expectedNatsServers:       defaultNatsServers,
			expectedNatsUser:          defaultNatsUser,
			expectedNatsPassword:      defaultNatsPassword,
			expectedCockroachAddr:     defaultCockroachAddr,
			expectedCockroachUser:     defaultCockroachUser,
			expectedCockroachPassword: defaultCockroachPassword,
			expectedCockroachDatabase: defaultCockroachDatabase,
			expectedJaegerAgentAddr:   defaultJaegerAgentAddr,
			expectedJaegerLogSpans:    defaultJaegerLogSpans,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedLogLevel, Config.LogLevel)
			assert.Equal(t, tc.expectedServiceName, Config.ServiceName)
			assert.Equal(t, tc.expectedServicePort, Config.ServicePort)
			assert.Equal(t, tc.expectedNatsServers, Config.NatsServers)
			assert.Equal(t, tc.expectedNatsUser, Config.NatsUser)
			assert.Equal(t, tc.expectedNatsPassword, Config.NatsPassword)
			assert.Equal(t, tc.expectedCockroachAddr, Config.CockroachAddr)
			assert.Equal(t, tc.expectedCockroachUser, Config.CockroachUser)
			assert.Equal(t, tc.expectedCockroachPassword, Config.CockroachPassword)
			assert.Equal(t, tc.expectedCockroachDatabase, Config.CockroachDatabase)
			assert.Equal(t, tc.expectedJaegerAgentAddr, Config.JaegerAgentAddr)
			assert.Equal(t, tc.expectedJaegerLogSpans, Config.JaegerLogSpans)
		})
	}
}
