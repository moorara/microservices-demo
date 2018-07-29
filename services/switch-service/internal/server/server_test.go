package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/internal/service"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

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
		name           string
		httpPort       string
		grpcPort       string
		caChainFile    string
		serverCertFile string
		serverKeyFile  string
		switchService  proto.SwitchServiceServer
	}{
		{
			name:          "MTLSDisabled",
			httpPort:      ":9998",
			grpcPort:      ":9999",
			switchService: &service.MockSwitchService{},
		},
		{
			name:           "MTLSEnabled",
			httpPort:       ":9998",
			grpcPort:       ":9999",
			caChainFile:    "../test/ca.chain.cert",
			serverCertFile: "../test/server.cert",
			serverKeyFile:  "../test/server.key",
			switchService:  &service.MockSwitchService{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			server, err := New(
				tc.httpPort, tc.grpcPort,
				tc.caChainFile, tc.serverCertFile, tc.serverKeyFile,
				tc.switchService, logger, metrics,
			)

			assert.NoError(t, err)
			assert.NotNil(t, server.logger)
			assert.Equal(t, tc.httpPort, server.httpPort)
			assert.NotNil(t, server.httpServer)
			assert.Equal(t, tc.grpcPort, server.grpcPort)
			assert.NotNil(t, server.grpcServer)
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

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestShutdown(t *testing.T) {
	tests := []struct {
		name       string
		httpServer *mockHTTPServer
		grpcServer *mockGRPCServer
	}{
		{
			"HTTPServerError",
			&mockHTTPServer{
				ShutdownOutError: errors.New("http error"),
			},
			&mockGRPCServer{},
		},
		{
			"NoError",
			&mockHTTPServer{},
			&mockGRPCServer{},
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

			server.Shutdown()

			assert.True(t, tc.httpServer.ShutdownCalled)
			assert.True(t, tc.grpcServer.GracefulStopCalled)
		})
	}
}
