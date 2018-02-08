package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		logLevel    string
		serviceName string
		loggerName  string
	}{
		{
			"Debug", "debug", "job-service", "go-kit",
		},
		{
			"Info", "info", "auth-service", "go-kit",
		},
		{
			"Warn", "warn", "gateway-service", "go-kit",
		},
		{
			"Error", "error", "storage-service", "go-kit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewLogger(tc.logLevel, tc.serviceName, tc.loggerName)

			assert.NotNil(t, logger)
		})
	}
}
