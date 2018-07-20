package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/sensor-service/util"
	"github.com/opentracing/opentracing-go/log"

	opentracing "github.com/opentracing/opentracing-go"
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

func (tm *tracerMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		url := r.URL.Path

		// This only works with mux router
		for p, v := range mux.Vars(r) {
			url = strings.Replace(url, v, ":"+p, 1)
		}

		rw := util.NewResponseWriter(w)
		next(rw, r)

		statusCode := rw.StatusCode()

		span := tm.tracer.StartSpan("http-request")
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
