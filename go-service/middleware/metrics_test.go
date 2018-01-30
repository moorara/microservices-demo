package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moorara/toys/microservices/go-service/util"
	"github.com/stretchr/testify/assert"
)

func TestGetMetricsWrapper(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"101", 101},
		{"200", 200},
		{"302", 302},
		{"404", 404},
		{"500", 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metrics := util.NewMetrics("go_service")
			metricsMiddleware := NewMetricsMiddleware(metrics)

			r := httptest.NewRequest("GET", "http://service/votes", nil)
			w := httptest.NewRecorder()

			handler := metricsMiddleware.Wrap(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})
			handler(w, r)
			res := w.Result()

			assert.Equal(t, tc.statusCode, res.StatusCode)
		})
	}
}
