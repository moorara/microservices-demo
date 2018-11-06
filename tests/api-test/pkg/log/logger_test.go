package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	LogPassedKV []interface{}
	LogError    error
}

func (m *mockLogger) Log(kv ...interface{}) error {
	m.LogPassedKV = kv
	return m.LogError
}

func TestNewNopLogger(t *testing.T) {
	logger := NewNopLogger()
	assert.NotNil(t, logger)
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name       string
		loggerName string
		logLevel   string
	}{
		{"LevelDebug", "api-test", "debug"},
		{"LevelInfo", "api-test", "info"},
		{"LevelWarn", "test-service", "warn"},
		{"LevelError", "test-service", "error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewLogger(tc.loggerName, tc.logLevel)
			assert.NotNil(t, logger)
		})
	}
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger mockLogger
		kv         []interface{}
		expectedKV []interface{}
	}{
		{
			"OK",
			mockLogger{},
			[]interface{}{"key", "value", "message", "content"},
			[]interface{}{"key", "value", "message", "content"},
		},
		{
			"Error",
			mockLogger{
				LogError: errors.New("error"),
			},
			[]interface{}{"key", "value", "message", "content"},
			[]interface{}{"key", "value", "message", "content"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("DebugLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Debug(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})

			t.Run("InfoLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Info(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})

			t.Run("WarnLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Warn(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})

			t.Run("ErrorLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Error(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})
		})
	}
}
