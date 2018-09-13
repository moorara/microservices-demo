package db

import (
	"testing"

	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

type mockKitLogger struct {
	LogCalled   bool
	LogInKV     []interface{}
	LogOutError error
}

func (m *mockKitLogger) Log(kv ...interface{}) error {
	m.LogInKV = kv
	return m.LogOutError
}

func TestGormLogger(t *testing.T) {
	tests := []struct {
		values []interface{}
	}{
		{
			[]interface{}{"test.go:27", "invalid input"},
		},
	}

	for _, tc := range tests {
		kitLogger := &mockKitLogger{}
		logger := &log.Logger{
			Logger: kitLogger,
		}

		glogger := &gormLogger{
			logger: logger,
		}

		glogger.Print(tc.values...)
		assert.Contains(t, kitLogger.LogInKV, tc.values[1])
	}
}

func TestNewCockroachORM(t *testing.T) {
	tests := []struct {
		name        string
		addr        string
		user        string
		password    string
		database    string
		expectError bool
	}{
		{"WithoutUser", "localhost:26257", "", "", "things", true},
		{"WithUser", "localhost:26257", "root", "", "things", true},
		{"WithUserPass", "localhost:26257", "service", "password", "things", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			orm, err := NewCockroachORM(tc.addr, tc.user, tc.password, tc.database, logger)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, orm)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, orm)
			}
		})
	}
}
