package server

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moorara/microservices-demo/services/sensor-service/config"
	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	ListenAndServeCalled bool
	ListenAndServeError  error

	ShutdownCalled bool
	ShutdownError  error
}

func (s *mockServer) ListenAndServe() error {
	s.ListenAndServeCalled = true
	return s.ListenAndServeError
}

func (s *mockServer) Shutdown(context.Context) error {
	s.ShutdownCalled = true
	return s.ShutdownError
}

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
	}{
		{
			"Server1",
			config.Config{
				LogLevel:               "info",
				ServiceName:            "go-service",
				ServicePort:            ":4020",
				PostgresURL:            "postgres://localhost",
				JaegerAgentHost:        "localhost",
				JaegerAgentPort:        6831,
				JaegerReporterLogSpans: false,
			},
		},
		{
			"Server2",
			config.Config{
				LogLevel:               "debug",
				ServiceName:            "sensor-service",
				ServicePort:            ":4020",
				PostgresURL:            "postgres://root:pass@localhost",
				JaegerAgentHost:        "localhost",
				JaegerAgentPort:        6831,
				JaegerReporterLogSpans: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New(tc.config)

			assert.Equal(t, tc.config, s.config)
			assert.NotNil(t, s.server)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name                string
		logger              log.Logger
		listenAndServeError error
	}{
		{
			"ServerError",
			log.NewNopLogger(),
			errors.New("error"),
		},
		{
			"Successful",
			log.NewNopLogger(),
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &HTTPServer{
				config: config.Config{
					ServicePort: ":4020",
				},
				logger: tc.logger,
				server: &mockServer{
					ListenAndServeError: tc.listenAndServeError,
				},
			}

			err := s.Start()
			assert.Equal(t, tc.listenAndServeError, err)
		})
	}
}
