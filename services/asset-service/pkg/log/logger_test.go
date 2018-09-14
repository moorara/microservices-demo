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
		name        string
		serviceName string
		loggerName  string
		logLevel    string
	}{
		{"LevelDebug", "go-service", "singleton", "debug"},
		{"LevelInfo", "go-service", "singleton", "info"},
		{"LevelWarn", "go-service", "singleton", "warn"},
		{"LevelError", "go-service", "singleton", "error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLogger(tc.serviceName, tc.loggerName, tc.logLevel)
			assert.NotNil(t, l)
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
				l := &Logger{Logger: &tc.mockLogger}
				err := l.Debug(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})

			t.Run("InfoLevel", func(t *testing.T) {
				l := &Logger{Logger: &tc.mockLogger}
				err := l.Info(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})

			t.Run("WarnLevel", func(t *testing.T) {
				l := &Logger{Logger: &tc.mockLogger}
				err := l.Warn(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})

			t.Run("ErrorLevel", func(t *testing.T) {
				l := &Logger{Logger: &tc.mockLogger}
				err := l.Error(tc.kv)
				assert.Equal(t, tc.mockLogger.LogError, err)
				assert.Equal(t, tc.expectedKV, tc.mockLogger.LogPassedKV[2])
			})
		})
	}
}
