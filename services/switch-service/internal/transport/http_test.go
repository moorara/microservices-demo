package transport

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPServer(t *testing.T) {
	tests := []struct {
		name                      string
		addr                      string
		liveHandler               http.HandlerFunc
		readyHandler              http.HandlerFunc
		metricsHandler            http.HandlerFunc
		expectedLiveStatusCode    int
		expectedReadyStatusCode   int
		expectedMetricsStatusCode int
	}{
		{
			"LiveNotReady",
			"127.0.0.1:12345",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			http.StatusOK,
			http.StatusAccepted,
			http.StatusNotFound,
		},
		{
			"LiveAndReady",
			"127.0.0.1:12345",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			http.StatusOK,
			http.StatusOK,
			http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpServer := NewHTTPServer(tc.addr, tc.liveHandler, tc.readyHandler, tc.metricsHandler)
			assert.NotNil(t, httpServer)

			// Start http server in a goroutine
			errs := make(chan error)
			go func() {
				errs <- httpServer.ListenAndServe()
			}()

			t.Run("Liveness", func(t *testing.T) {
				resp, err := http.Get("http://" + tc.addr + "/live")
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLiveStatusCode, resp.StatusCode)
			})

			t.Run("Readiness", func(t *testing.T) {
				resp, err := http.Get("http://" + tc.addr + "/ready")
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedReadyStatusCode, resp.StatusCode)
			})

			t.Run("Metrics", func(t *testing.T) {
				resp, err := http.Get("http://" + tc.addr + "/metrics")
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMetricsStatusCode, resp.StatusCode)
			})

			httpServer.Close()
			err := <-errs
			assert.Equal(t, http.ErrServerClosed, err)
		})
	}
}
