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
			assert.Equal(t, tc.expectedLogLevel, Global.LogLevel)
			assert.Equal(t, tc.expectedServiceName, Global.ServiceName)
			assert.Equal(t, tc.expectedServicePort, Global.ServicePort)
			assert.Equal(t, tc.expectedNatsServers, Global.NatsServers)
			assert.Equal(t, tc.expectedNatsUser, Global.NatsUser)
			assert.Equal(t, tc.expectedNatsPassword, Global.NatsPassword)
			assert.Equal(t, tc.expectedCockroachAddr, Global.CockroachAddr)
			assert.Equal(t, tc.expectedCockroachUser, Global.CockroachUser)
			assert.Equal(t, tc.expectedCockroachPassword, Global.CockroachPassword)
			assert.Equal(t, tc.expectedCockroachDatabase, Global.CockroachDatabase)
			assert.Equal(t, tc.expectedJaegerAgentAddr, Global.JaegerAgentAddr)
			assert.Equal(t, tc.expectedJaegerLogSpans, Global.JaegerLogSpans)
		})
	}
}
