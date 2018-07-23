package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go"
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
		{"101", "GET", "http://service/resource/1111", 101, "/resource/1111"},
		{"200", "GET", "http://service/resource/2222", 200, "/resource/2222"},
		{"302", "GET", "http://service/resource/3333", 302, "/resource/3333"},
		{"404", "GET", "http://service/resource/4444", 404, "/resource/4444"},
		{"500", "GET", "http://service/resource/5555", 500, "/resource/5555"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var spanFromCtx opentracing.Span
			tracer := mocktracer.New()
			tracerMiddleware := NewTracerMiddleware(tracer)

			r := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			handler := tracerMiddleware.Wrap(func(w http.ResponseWriter, r *http.Request) {
				spanFromCtx = opentracing.SpanFromContext(r.Context())
				w.WriteHeader(tc.statusCode)
			})
			handler(w, r)

			spans := tracer.FinishedSpans()

			assert.Equal(t, spanFromCtx, spans[0])
			assert.Equal(t, "http-request", spans[0].OperationName)
			assert.Equal(t, tc.method, spans[0].Tag("http.method"))
			assert.Equal(t, tc.expectedURL, spans[0].Tag("http.url"))
			assert.Equal(t, tc.statusCode, spans[0].Tag("http.status_code"))
		})
	}
}
