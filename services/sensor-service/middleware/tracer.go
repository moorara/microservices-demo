package middleware

import (
	"net/http"

	"github.com/moorara/microservices-demo/services/sensor-service/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type tracerMiddleware struct {
	tracer opentracing.Tracer
}

// NewTracerMiddleware creates a new middleware for tracing
func NewTracerMiddleware(tracer opentracing.Tracer) Middleware {
	return &tracerMiddleware{
		tracer: tracer,
	}
}

func (m *tracerMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := m.tracer.StartSpan("http-request")
		ctx := opentracing.ContextWithSpan(r.Context(), span)

		rw := util.NewResponseWriter(w)
		req := r.WithContext(ctx)
		next(rw, req)

		method := r.Method
		url := r.URL.Path
		statusCode := rw.StatusCode()

		// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
		span.SetTag("http.method", method)
		span.SetTag("http.url", url)
		span.SetTag("http.status_code", statusCode)
		span.LogFields(
			log.String("event", "http-request"),
			log.String("method", method),
			log.String("url", url),
			log.Int("statusCode", statusCode),
		)

		span.Finish()
	}
}
