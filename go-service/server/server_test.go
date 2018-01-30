package server

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moorara/toys/microservices/go-service/config"
	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	listenAndServeCalled bool
	listenAndServeError  error

	shutdownCalled bool
	shutdownError  error
}

func (s *mockServer) ListenAndServe() error {
	s.listenAndServeCalled = true
	return s.listenAndServeError
}

func (s *mockServer) Shutdown(context.Context) error {
	s.shutdownCalled = true
	return s.shutdownError
}

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
	}{
		{
			"Server1",
			config.Config{
				LogLevel:    "info",
				ServiceName: "go-service",
				ServicePort: ":4010",
				RedisURL:    "redis://redis:6379",
			},
		},
		{
			"Server2",
			config.Config{
				LogLevel:    "debug",
				ServiceName: "golang-service",
				ServicePort: ":4020",
				RedisURL:    "redis://user:pass@redis:6389",
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
					ServicePort: ":4010",
				},
				logger: tc.logger,
				server: &mockServer{
					listenAndServeError: tc.listenAndServeError,
				},
			}

			err := s.Start()
			assert.Equal(t, tc.listenAndServeError, err)
		})
	}
}
