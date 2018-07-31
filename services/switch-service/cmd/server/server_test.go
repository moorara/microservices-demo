package server

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		httpPort      string
		grpcPort      string
		caFile        string
		certFile      string
		keyFile       string
		switchService proto.SwitchServiceServer
		expectError   bool
	}{
		{
			name:          "InvalidTLS",
			httpPort:      ":12345",
			grpcPort:      ":12346",
			caFile:        "ca.chain.cert",
			certFile:      "server.cert",
			keyFile:       "server.key",
			switchService: &mockSwitchService{},
			expectError:   true,
		},
		{
			name:          "Simple",
			httpPort:      ":12345",
			grpcPort:      ":12346",
			switchService: &mockSwitchService{},
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			server, err := New(
				tc.httpPort, tc.grpcPort,
				tc.caFile, tc.certFile, tc.keyFile,
				tc.switchService, logger, metrics,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
			}
		})
	}
}

func TestLiveHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/live", nil)
	w := httptest.NewRecorder()

	server := &server{}
	server.liveHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestReadyHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	server := &server{}
	server.readyHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestStart(t *testing.T) {
	tests := []struct {
		name              string
		grpcPort          string
		httpServer        *mockHTTPServer
		grpcServer        *mockGRPCServer
		signal            syscall.Signal
		expectedErrorType error
	}{
		{
			"ReservedPort",
			":80",
			&mockHTTPServer{},
			&mockGRPCServer{},
			syscall.SIGINT,
			&net.OpError{},
		},
		{
			"InterruptSignal",
			"", // Random port
			&mockHTTPServer{},
			&mockGRPCServer{},
			syscall.SIGINT,
			errors.New("interrupt"),
		},
		{
			"TerminationSignal",
			"", // Random port
			&mockHTTPServer{},
			&mockGRPCServer{},
			syscall.SIGTERM,
			errors.New("terminated"),
		},
		{
			"HTTPServerError",
			"", // Random port
			&mockHTTPServer{
				ListenAndServeOutError: errors.New("http error"),
			},
			&mockGRPCServer{},
			0,
			errors.New("http error"),
		},
		{
			"GRPCServerError",
			"", // Random port
			&mockHTTPServer{},
			&mockGRPCServer{
				ServeOutError: errors.New("grpc error"),
			},
			0,
			errors.New("grpc error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			server := &server{
				logger:     logger,
				grpcPort:   tc.grpcPort,
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

			assert.IsType(t, tc.expectedErrorType, err)
		})
	}
}

func TestStop(t *testing.T) {
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
			server := &server{
				logger:     logger,
				httpServer: tc.httpServer,
				grpcServer: tc.grpcServer,
			}

			server.Stop()

			assert.Equal(t, 1, tc.httpServer.ShutdownCallCount)
			assert.Equal(t, 1, tc.grpcServer.GracefulStopCallCount)
		})
	}
}
