package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	xhttp "github.com/moorara/microservices-demo/services/asset-service/pkg/http"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	spanName = "http-request"
)

type monitorMiddleware struct {
	logger  *log.Logger
	metrics *metrics.Metrics
	tracer  opentracing.Tracer
}

// NewMonitorMiddleware creates a new middleware for logging
func NewMonitorMiddleware(logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) Middleware {
	return &monitorMiddleware{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (m *monitorMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new trace
		span := m.tracer.StartSpan(spanName)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(r.Context(), span)

		// Next http handler
		start := time.Now()
		rw := xhttp.NewResponseWriter(w)
		req := r.WithContext(ctx)
		next(rw, req)
		duration := time.Now().Sub(start).Seconds()

		method := r.Method
		url := r.URL.Path
		headers := r.Header
		statusCode := uint16(rw.StatusCode())
		statusClass := rw.StatusClass()

		logs := []interface{}{
			"req.method", method,
			"req.url", url,
			"req.headers", headers,
			"res.statusCode", statusCode,
			"res.statusClass", statusClass,
			"responseTime", duration,
			"message", fmt.Sprintf("%s %s %d %f", method, url, statusCode, duration),
		}

		// Logging
		switch {
		case statusCode >= 500:
			m.logger.Error(logs...)
		case statusCode >= 400:
			m.logger.Warn(logs...)
		case statusCode >= 100:
			m.logger.Info(logs...)
		}

		// Metrics
		sc := strconv.Itoa(rw.StatusCode())
		m.metrics.HTTPDurationHist.WithLabelValues(method, url, sc, statusClass).Observe(duration)
		m.metrics.HTTPDurationSumm.WithLabelValues(method, url, sc, statusClass).Observe(duration)

		// Tracing
		// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
		ext.HTTPMethod.Set(span, method)
		ext.HTTPUrl.Set(span, url)
		ext.HTTPStatusCode.Set(span, statusCode)
	}
}
