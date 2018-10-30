package server

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moorara/microservices-demo/services/sensor-service/config"
	"github.com/stretchr/testify/assert"
)

type mockCloser struct {
	CloseCalled bool
	CloseError  error
}

func (c *mockCloser) Close() error {
	c.CloseCalled = true
	return c.CloseError
}

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
				LogLevel:         "info",
				ServiceName:      "go-service",
				ServicePort:      ":4020",
				PostgresHost:     "localhost",
				PostgresPort:     "5432",
				PostgresDatabase: "sensors",
				JaegerAgentAddr:  "localhost:6831",
				JaegerLogSpans:   false,
			},
		},
		{
			"Server2",
			config.Config{
				LogLevel:         "debug",
				ServiceName:      "sensor-service",
				ServicePort:      ":4020",
				PostgresHost:     "localhost",
				PostgresPort:     "5432",
				PostgresDatabase: "sensors",
				PostgresUsername: "root",
				PostgresPassword: "pass",
				JaegerAgentAddr:  "localhost:6831",
				JaegerLogSpans:   true,
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
			s := HTTPServer{
				config: config.Config{
					ServicePort: ":4020",
				},
				logger:  tc.logger,
				closers: []io.Closer{},
				server: &mockServer{
					ListenAndServeError: tc.listenAndServeError,
				},
			}

			err := s.Start()
			assert.Equal(t, tc.listenAndServeError, err)
		})
	}
}

func TestClose(t *testing.T) {
	tests := []struct {
		name        string
		closers     []io.Closer
		expectError bool
	}{
		{
			"NoCloser",
			[]io.Closer{},
			false,
		},
		{
			"Error",
			[]io.Closer{
				&mockCloser{},
				&mockCloser{CloseError: errors.New("error")},
				&mockCloser{},
			},
			true,
		},
		{
			"NoError",
			[]io.Closer{
				&mockCloser{},
				&mockCloser{},
				&mockCloser{},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := HTTPServer{
				closers: tc.closers,
			}

			err := s.Close()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
