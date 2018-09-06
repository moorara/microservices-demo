package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	kitLog "github.com/go-kit/kit/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func getMetric(metrics *metrics.Metrics, name string) (output string) {
	mfs, _ := metrics.Registry.Gather()
	for _, mf := range mfs {
		if *mf.Name == name {
			for _, m := range mf.Metric {
				output += m.String() + "\n"
			}
		}
	}
	return output
}

func TestMonitorMiddleware(t *testing.T) {
	tests := []struct {
		name                string
		parentSpan          opentracing.Span
		reqMethod           string
		reqURL              string
		resStatusCode       int
		expectedMethod      string
		expectedURL         string
		expectedStatusCode  int
		expectedStatusClass string
	}{
		{
			name:                "200",
			parentSpan:          mocktracer.New().StartSpan("parent-span"),
			reqMethod:           "POST",
			reqURL:              "/graphql",
			resStatusCode:       200,
			expectedMethod:      "POST",
			expectedURL:         "/graphql",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
		{
			name:                "404",
			parentSpan:          nil,
			reqMethod:           "POST",
			reqURL:              "/graphql",
			resStatusCode:       404,
			expectedMethod:      "POST",
			expectedURL:         "/graphql",
			expectedStatusCode:  404,
			expectedStatusClass: "4xx",
		},
		{
			name:                "500",
			parentSpan:          nil,
			reqMethod:           "POST",
			reqURL:              "/graphql",
			resStatusCode:       500,
			expectedMethod:      "POST",
			expectedURL:         "/graphql",
			expectedStatusCode:  500,
			expectedStatusClass: "5xx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var spanFromCtx opentracing.Span

			// Logger with pipe to read from
			rd, wr, _ := os.Pipe()
			dec := json.NewDecoder(rd)
			logger := &log.Logger{
				Logger: kitLog.NewJSONLogger(wr),
			}

			// Mock metrics and tracer
			metrics := metrics.New("test_service")
			tracer := mocktracer.New()

			middleware := NewMonitorMiddleware(logger, metrics, tracer)
			handler := middleware.Wrap(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resStatusCode)
				spanFromCtx = opentracing.SpanFromContext(r.Context())
			})

			r := httptest.NewRequest(tc.reqMethod, tc.reqURL, nil)
			w := httptest.NewRecorder()

			// Inject parent span context
			if tc.parentSpan != nil {
				carrier := opentracing.HTTPHeadersCarrier(r.Header)
				err := tracer.Inject(tc.parentSpan.Context(), opentracing.HTTPHeaders, carrier)
				assert.NoError(t, err)
			}

			handler(w, r)

			// Verify logging
			var log map[string]interface{}
			dec.Decode(&log)
			assert.Equal(t, tc.expectedMethod, log["req.method"])
			assert.Equal(t, tc.expectedURL, log["req.url"])
			assert.Equal(t, float64(tc.expectedStatusCode), log["res.statusCode"])
			assert.Equal(t, tc.expectedStatusClass, log["res.statusClass"])

			// Verify metrics
			assert.NotEmpty(t, getMetric(metrics, "http_requests_duration_seconds"))
			assert.NotEmpty(t, getMetric(metrics, "http_requests_duration_quantiles_seconds"))

			// Verify tracing
			span := tracer.FinishedSpans()[0]
			assert.Equal(t, spanFromCtx, span)
			assert.Equal(t, "http-request", span.OperationName)
			assert.Equal(t, tc.expectedMethod, span.Tag("http.method"))
			assert.Equal(t, tc.expectedURL, span.Tag("http.url"))
			assert.Equal(t, uint16(tc.expectedStatusCode), span.Tag("http.status_code"))

			if tc.parentSpan != nil {
				parentSpan, _ := tc.parentSpan.(*mocktracer.MockSpan)
				assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
				assert.Equal(t, parentSpan.SpanContext.TraceID, span.SpanContext.TraceID)
			}
		})
	}
}
