package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/moorara/toys/microservices/go-service/util"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	histogramName    = "http_requests_duration_seconds"
	summaryName      = "http_requests_duration_quantiles_seconds"
	defaultLabels    = []string{"method", "endpoint", "statusCode", "statusClass"}
	defaultBuckets   = []float64{0.01, 0.1, 0.5, 1.}
	defaultQuantiles = map[float64]float64{0.1: 0.1, 0.5: 0.05, 0.95: 0.01, 0.99: 0.001}
)

type metricsMiddleware struct {
	histogram *prometheus.HistogramVec
	summary   *prometheus.SummaryVec
}

// NewMetricsMiddleware creates a new middleware for metrics
func NewMetricsMiddleware(metrics *util.Metrics) Middleware {
	return &metricsMiddleware{
		histogram: metrics.NewHistogram(
			histogramName,
			"duration histogram of http requests",
			defaultBuckets,
			defaultLabels,
		),
		summary: metrics.NewSummary(
			summaryName,
			"duration summary of http requests",
			defaultQuantiles,
			defaultLabels,
		),
	}
}

func (mm *metricsMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		endpoint := r.URL.Path

		// This only works with mux router
		for p, v := range mux.Vars(r) {
			endpoint = strings.Replace(endpoint, v, ":"+p, 1)
		}

		rw := util.NewResponseWriter(w)
		next(rw, r)

		duration := time.Now().Sub(start).Seconds()
		statusCode := strconv.Itoa(rw.StatusCode())
		statusClass := rw.StatusClass()

		mm.histogram.WithLabelValues(method, endpoint, statusCode, statusClass).Observe(duration)
		mm.summary.WithLabelValues(method, endpoint, statusCode, statusClass).Observe(duration)
	}
}
