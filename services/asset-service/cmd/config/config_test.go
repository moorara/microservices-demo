package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name                    string
		expectedLogLevel        string
		expectedServiceName     string
		expectedServicePort     string
		expectedArangoEndpoints []string
		expectedArangoUser      string
		expectedArangoPassword  string
		expectedNatsServers     []string
		expectedNatsUser        string
		expectedNatsPassword    string
		expectedJaegerAgentAddr string
		expectedJaegerLogSpans  bool
	}{
		{
			name:                    "Defauts",
			expectedLogLevel:        defaultLogLevel,
			expectedServiceName:     defaultServiceName,
			expectedServicePort:     defaultServicePort,
			expectedArangoEndpoints: defaultArangoEndpoints,
			expectedArangoUser:      defaultArangoUser,
			expectedArangoPassword:  defaultArangoPassword,
			expectedNatsServers:     defaultNatsServers,
			expectedNatsUser:        defaultNatsUser,
			expectedNatsPassword:    defaultNatsPassword,
			expectedJaegerAgentAddr: defaultJaegerAgentAddr,
			expectedJaegerLogSpans:  defaultJaegerLogSpans,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedLogLevel, Config.LogLevel)
			assert.Equal(t, tc.expectedServiceName, Config.ServiceName)
			assert.Equal(t, tc.expectedServicePort, Config.ServicePort)
			assert.Equal(t, tc.expectedArangoEndpoints, Config.ArangoEndpoints)
			assert.Equal(t, tc.expectedArangoUser, Config.ArangoUser)
			assert.Equal(t, tc.expectedArangoPassword, Config.ArangoPassword)
			assert.Equal(t, tc.expectedNatsServers, Config.NatsServers)
			assert.Equal(t, tc.expectedNatsUser, Config.NatsUser)
			assert.Equal(t, tc.expectedNatsPassword, Config.NatsPassword)
			assert.Equal(t, tc.expectedJaegerAgentAddr, Config.JaegerAgentAddr)
			assert.Equal(t, tc.expectedJaegerLogSpans, Config.JaegerLogSpans)
		})
	}
}
