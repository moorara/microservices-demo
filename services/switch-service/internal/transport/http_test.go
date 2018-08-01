package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPServer(t *testing.T) {
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
			"http://localhost:9999",
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
			"http://localhost:9999",
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

			server, ok := httpServer.(*http.Server)
			assert.True(t, ok)

			ts := httptest.NewServer(server.Handler)
			defer ts.Close()

			t.Run("Liveness", func(t *testing.T) {
				resp, err := http.Get(ts.URL + "/live")
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLiveStatusCode, resp.StatusCode)
			})

			t.Run("Readiness", func(t *testing.T) {
				resp, err := http.Get(ts.URL + "/ready")
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedReadyStatusCode, resp.StatusCode)
			})

			t.Run("Metrics", func(t *testing.T) {
				resp, err := http.Get(ts.URL + "/metrics")
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMetricsStatusCode, resp.StatusCode)
			})
		})
	}
}
