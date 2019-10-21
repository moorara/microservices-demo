package middleware

import (
	"net/http"

	"github.com/moorara/microservices-demo/services/sensor/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const spanName = "http-request"

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
		var span opentracing.Span

		carrier := opentracing.HTTPHeadersCarrier(r.Header)
		parentSpanContext, err := m.tracer.Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			span = m.tracer.StartSpan(spanName)
		} else {
			span = m.tracer.StartSpan(spanName, opentracing.ChildOf(parentSpanContext))
		}

		ctx := opentracing.ContextWithSpan(r.Context(), span)
		rw := util.NewResponseWriter(w)
		req := r.WithContext(ctx)
		next(rw, req)

		proto := r.Proto
		method := r.Method
		url := r.URL.Path
		statusCode := uint16(rw.StatusCode())

		// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
		span.SetTag("http.version", proto)
		ext.HTTPMethod.Set(span, method)
		ext.HTTPUrl.Set(span, url)
		ext.HTTPStatusCode.Set(span, statusCode)
		span.Finish()
	}
}
