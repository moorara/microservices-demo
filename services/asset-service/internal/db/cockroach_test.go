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
		name     string
		addr     string
		user     string
		password string
		database string
	}{
		{"WithoutUser", "localhost:26257", "", "", "things"},
		{"WithUser", "localhost:26257", "root", "", "things"},
		{"WithUserPass", "localhost:26257", "service", "password", "things"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Log("recovered from panic!")
				}
			}()

			logger := log.NewNopLogger()
			db := NewCockroachORM(tc.addr, tc.user, tc.password, tc.database, logger)
			assert.NotNil(t, db)
		})
	}
}
