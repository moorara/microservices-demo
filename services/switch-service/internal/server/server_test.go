package server

import (
	"errors"
	"io"
	"net"
	"testing"

	"github.com/moorara/microservices-demo/services/switch-service/cmd/config"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

type mockCloser struct {
	CloseCalled   bool
	CloseOutError error
}

func (m *mockCloser) Close() error {
	m.CloseCalled = true
	return m.CloseOutError
}

type mockGRPCServer struct {
	ServeCalled        bool
	ServeInListener    net.Listener
	ServeOutError      error
	GracefulStopCalled bool
}

func (m *mockGRPCServer) Serve(listener net.Listener) error {
	m.ServeCalled = true
	m.ServeInListener = listener
	return m.ServeOutError
}

func (m *mockGRPCServer) GracefulStop() {
	m.GracefulStopCalled = true
}

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
	}{
		{
			"MTLSDisabled",
			config.Config{
				ServiceName: "go-service",
				ServicePort: ":9999",
			},
		},
		{
			"MTLSEnabled",
			config.Config{
				ServiceName:    "go-service",
				ServicePort:    ":9999",
				CAChainFile:    "../test/ca.chain.cert",
				ServerCertFile: "../test/server.cert",
				ServerKeyFile:  "../test/server.key",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server, err := New(tc.config)

			assert.NoError(t, err)
			assert.Equal(t, tc.config.ServicePort, server.port)
			assert.NotNil(t, server.logger)
			assert.NotNil(t, server.grpcServer)
			assert.NotEmpty(t, server.closers)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name       string
		grpcServer *mockGRPCServer
	}{
		{
			"GRPCServerError",
			&mockGRPCServer{
				ServeOutError: errors.New("error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			server := &Server{
				logger:     logger,
				grpcServer: tc.grpcServer,
			}

			err := server.Start()
			defer tc.grpcServer.ServeInListener.Close()

			assert.True(t, tc.grpcServer.ServeCalled)
			assert.Equal(t, tc.grpcServer.ServeOutError, err)
		})
	}
}

func TestShutdown(t *testing.T) {
	tests := []struct {
		name        string
		grpcServer  *mockGRPCServer
		closers     []io.Closer
		expectError bool
	}{
		{
			"NoCloser",
			&mockGRPCServer{},
			[]io.Closer{},
			false,
		},
		{
			"NoError",
			&mockGRPCServer{},
			[]io.Closer{
				&mockCloser{},
				&mockCloser{},
			},
			false,
		},
		{
			"CloserError",
			&mockGRPCServer{},
			[]io.Closer{
				&mockCloser{},
				&mockCloser{CloseOutError: errors.New("error")},
			},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			server := &Server{
				logger:     logger,
				grpcServer: tc.grpcServer,
				closers:    tc.closers,
			}

			err := server.Shutdown()

			assert.True(t, tc.grpcServer.GracefulStopCalled)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
