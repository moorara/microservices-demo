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
		name               string
		reqMethod          string
		reqURL             string
		parentSpan         opentracing.Span
		resStatusCode      int
		expectedVersion    string
		expectedMethod     string
		expectedURL        string
		expectedStatusCode uint16
	}{
		{
			name:               "WithParentSpan",
			reqMethod:          "GET",
			reqURL:             "http://service/resource/1111",
			parentSpan:         mocktracer.New().StartSpan("parent-span"),
			resStatusCode:      200,
			expectedVersion:    "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/resource/1111",
			expectedStatusCode: 200,
		},
		{
			name:               "WithoutParentSpan",
			reqMethod:          "GET",
			reqURL:             "http://service/resource/4444",
			parentSpan:         nil,
			resStatusCode:      404,
			expectedVersion:    "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/resource/4444",
			expectedStatusCode: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var spanFromCtx opentracing.Span

			tracer := mocktracer.New()
			tracerMiddleware := NewTracerMiddleware(tracer)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.reqMethod, tc.reqURL, nil)

			// Inject parent span context
			if tc.parentSpan != nil {
				carrier := opentracing.HTTPHeadersCarrier(r.Header)
				err := tracer.Inject(tc.parentSpan.Context(), opentracing.HTTPHeaders, carrier)
				assert.NoError(t, err)
			}

			handler := tracerMiddleware.Wrap(func(w http.ResponseWriter, r *http.Request) {
				spanFromCtx = opentracing.SpanFromContext(r.Context())
				w.WriteHeader(tc.resStatusCode)
			})
			handler(w, r)

			span := tracer.FinishedSpans()[0]

			assert.Equal(t, spanFromCtx, span)
			assert.Equal(t, "http-request", span.OperationName)
			assert.Equal(t, tc.expectedVersion, span.Tag("http.version"))
			assert.Equal(t, tc.expectedMethod, span.Tag("http.method"))
			assert.Equal(t, tc.expectedURL, span.Tag("http.url"))
			assert.Equal(t, tc.expectedStatusCode, span.Tag("http.status_code"))

			if tc.parentSpan != nil {
				parentSpan, _ := tc.parentSpan.(*mocktracer.MockSpan)
				assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
				assert.Equal(t, parentSpan.SpanContext.TraceID, span.SpanContext.TraceID)
			}
		})
	}
}
