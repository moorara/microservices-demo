package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/switch/cmd/config"
	"github.com/moorara/microservices-demo/services/switch/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch/pkg/log"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		config      config.Config
		expectError bool
	}{
		{
			"InvalidTLS",
			config.Config{
				ServiceHTTPPort: ":12345",
				ServiceGRPCPort: ":12346",
				ArangoEndpoints: []string{"localhost:12347"},
				ArangoUser:      "root",
				ArangoPassword:  "pass",
				CAChainFile:     "ca.chain.cert",
				ServerCertFile:  "server.cert",
				ServerKeyFile:   "server.key",
			},
			true,
		},
		{
			"Simple",
			config.Config{
				ServiceHTTPPort: ":12345",
				ServiceGRPCPort: ":12346",
				ArangoEndpoints: []string{"localhost:12347"},
				ArangoUser:      "root",
				ArangoPassword:  "pass",
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			server, err := New(tc.config, logger, metrics, tracer)

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

			server := &server{
				ready: tc.ready,
			}

			server.readyHandler(w, r)
			assert.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name              string
		config            config.Config
		arangoService     *mockArangoService
		httpServer        *mockHTTPServer
		grpcServer        *mockGRPCServer
		signal            syscall.Signal
		expectedErrorType error
	}{
		{
			"InterruptSignal",
			config.Config{
				ServiceGRPCPort:  "",
				ServerTimeout:    int64(100 * time.Millisecond),
				ArangoDatabase:   "database",
				ArangoCollection: "collection",
			},
			&mockArangoService{},
			&mockHTTPServer{},
			&mockGRPCServer{},
			syscall.SIGINT,
			errors.New("interrupt"),
		},
		{
			"TerminationSignal",
			config.Config{
				ServiceGRPCPort:  "",
				ServerTimeout:    int64(100 * time.Millisecond),
				ArangoDatabase:   "database",
				ArangoCollection: "collection",
			},
			&mockArangoService{
				ConnectOutError: errors.New("database error"),
			},
			&mockHTTPServer{},
			&mockGRPCServer{},
			syscall.SIGTERM,
			errors.New("terminated"),
		},
		{
			"ArangoServiceConnectError",
			config.Config{
				ServiceGRPCPort:  "",
				ServerTimeout:    int64(100 * time.Millisecond),
				ArangoDatabase:   "database",
				ArangoCollection: "collection",
			},
			&mockArangoService{
				ConnectOutError: errors.New("database error"),
			},
			&mockHTTPServer{},
			&mockGRPCServer{},
			0,
			errors.New("database error"),
		},
		/* {
			"GRPCPortAccessDenied",
			config.Config{
				ServiceGRPCPort:  ":80",
				ServerTimeout:    int64(100 * time.Millisecond),
				ArangoDatabase:   "database",
				ArangoCollection: "collection",
			},
			&mockArangoService{},
			&mockHTTPServer{},
			&mockGRPCServer{},
			0,
			&net.OpError{},
		}, */
		{
			"HTTPServerError",
			config.Config{
				ServiceGRPCPort:  "",
				ServerTimeout:    int64(100 * time.Millisecond),
				ArangoDatabase:   "database",
				ArangoCollection: "collection",
			},
			&mockArangoService{},
			&mockHTTPServer{
				ListenAndServeOutError: errors.New("http error"),
			},
			&mockGRPCServer{},
			0,
			errors.New("http error"),
		},
		{
			"GRPCServerError",
			config.Config{
				ServiceGRPCPort:  "",
				ServerTimeout:    int64(100 * time.Millisecond),
				ArangoDatabase:   "database",
				ArangoCollection: "collection",
			},
			&mockArangoService{},
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
				config:        tc.config,
				logger:        logger,
				arangoService: tc.arangoService,
				httpServer:    tc.httpServer,
				grpcServer:    tc.grpcServer,
			}

			if tc.signal > 0 {
				sig := tc.signal // to prevent data race
				go func() {
					time.Sleep(50 * time.Millisecond)
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

			assert.True(t, tc.httpServer.ShutdownCalled)
			assert.True(t, tc.grpcServer.GracefulStopCalled)
		})
	}
}
