package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestWrapWithTracer(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		statusCode  int
		expectedURL string
	}{
		{"101", "GET", "http://service/resource", 101, "/resource"},
		{"200", "GET", "http://service/resource", 200, "/resource"},
		{"302", "GET", "http://service/resource", 302, "/resource"},
		{"404", "GET", "http://service/resource", 404, "/resource"},
		{"500", "GET", "http://service/resource", 500, "/resource"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer := mocktracer.New()
			tracerMiddleware := NewTracerMiddleware(tracer)

			r := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			handler := tracerMiddleware.Wrap(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})
			handler(w, r)

			res := w.Result()
			spans := tracer.FinishedSpans()

			assert.Equal(t, tc.statusCode, res.StatusCode)
			assert.Equal(t, "http-request", spans[0].OperationName)
			assert.Equal(t, tc.method, spans[0].Tag("http.method"))
			assert.Equal(t, tc.expectedURL, spans[0].Tag("http.url"))
			assert.Equal(t, tc.statusCode, spans[0].Tag("http.status_code"))
		})
	}
}
