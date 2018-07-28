package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

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

type mockHTTPServer struct {
	ListenAndServeCalled   bool
	ListenAndServeOutError error
	ShutdownCalled         bool
	ShutdownInContext      context.Context
	ShutdownOutError       error
}

func (m *mockHTTPServer) ListenAndServe() error {
	m.ListenAndServeCalled = true
	return m.ListenAndServeOutError
}

func (m *mockHTTPServer) Shutdown(ctx context.Context) error {
	m.ShutdownCalled = true
	m.ShutdownInContext = ctx
	return m.ShutdownOutError
}

type mockGRPCServer struct {
	ServeCalled        bool
	ServeInListener    net.Listener
	ServeOutError      error
	GracefulStopCalled bool
}

func (m *mockGRPCServer) Serve(listener net.Listener) error {
	if listener != nil {
		listener.Close()
	}

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
				ServiceName:     "go-service",
				ServiceGRPCPort: ":9998",
				ServiceHTTPPort: ":9999",
			},
		},
		{
			"MTLSEnabled",
			config.Config{
				ServiceName:     "go-service",
				ServiceGRPCPort: ":9998",
				ServiceHTTPPort: ":9999",
				CAChainFile:     "../test/ca.chain.cert",
				ServerCertFile:  "../test/server.cert",
				ServerKeyFile:   "../test/server.key",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server, err := New(tc.config)

			assert.NoError(t, err)
			assert.Equal(t, tc.config.ServiceGRPCPort, server.grpcPort)
			assert.NotNil(t, server.logger)
			assert.NotNil(t, server.grpcServer)
			assert.NotEmpty(t, server.closers)
		})
	}
}

func TestLiveHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/live", nil)
	w := httptest.NewRecorder()

	server := &Server{}
	server.LiveHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name               string
		ready              bool
		expectedStatusCode int
	}{
		{"Ready", true, http.StatusOK},
		{"NotReady", false, http.StatusAccepted},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/ready", nil)
			w := httptest.NewRecorder()

			server := &Server{ready: tc.ready}
			server.ReadyHandler(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name          string
		signal        syscall.Signal
		httpServer    *mockHTTPServer
		grpcServer    *mockGRPCServer
		expectedError error
	}{
		{
			"InterruptSignal",
			syscall.SIGINT,
			&mockHTTPServer{},
			&mockGRPCServer{},
			errors.New("interrupt"),
		},
		{
			"TerminationSignal",
			syscall.SIGTERM,
			&mockHTTPServer{},
			&mockGRPCServer{},
			errors.New("terminated"),
		},
		{
			"HTTPServerError",
			0,
			&mockHTTPServer{
				ListenAndServeOutError: errors.New("http error"),
			},
			&mockGRPCServer{},
			errors.New("http error"),
		},
		{
			"GRPCServerError",
			0,
			&mockHTTPServer{},
			&mockGRPCServer{
				ServeOutError: errors.New("grpc error"),
			},
			errors.New("grpc error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			server := &Server{
				logger:     logger,
				httpServer: tc.httpServer,
				grpcServer: tc.grpcServer,
			}

			if tc.signal > 0 {
				sig := tc.signal // to prevent data race
				go func() {
					time.Sleep(100 * time.Millisecond)
					syscall.Kill(syscall.Getpid(), sig)
				}()
			}

			err := server.Start()

			fmt.Println(tc.signal)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestShutdown(t *testing.T) {
	tests := []struct {
		name       string
		httpServer *mockHTTPServer
		grpcServer *mockGRPCServer
		closers    []io.Closer
	}{
		{
			"HTTPServerError",
			&mockHTTPServer{
				ShutdownOutError: errors.New("http error"),
			},
			&mockGRPCServer{},
			[]io.Closer{},
		},
		{
			"CloserError",
			&mockHTTPServer{},
			&mockGRPCServer{},
			[]io.Closer{
				&mockCloser{},
				&mockCloser{
					CloseOutError: errors.New("error"),
				},
			},
		},
		{
			"NoCloser",
			&mockHTTPServer{},
			&mockGRPCServer{},
			[]io.Closer{},
		},
		{
			"NoError",
			&mockHTTPServer{},
			&mockGRPCServer{},
			[]io.Closer{
				&mockCloser{},
				&mockCloser{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			server := &Server{
				logger:     logger,
				httpServer: tc.httpServer,
				grpcServer: tc.grpcServer,
				closers:    tc.closers,
			}

			server.Shutdown()

			assert.True(t, tc.httpServer.ShutdownCalled)
			assert.True(t, tc.grpcServer.GracefulStopCalled)
		})
	}
}
