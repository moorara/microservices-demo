package metrics

import (
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics defines all the metrics
type Metrics struct {
	Registry      *prometheus.Registry
	ReqCounter    *prometheus.CounterVec
	OpLatencyHist *prometheus.HistogramVec
	OpLatencySumm *prometheus.SummaryVec
}

// New creates a new metrics
func New(service string) *Metrics {
	service = strings.Replace(service, "-", "_", -1)

	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), service))

	ReqCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: service,
			Name:      "graphql_requests_total",
			Help:      "total number of graphql requests",
		},
		[]string{"call", "success"},
	)

	OpLatencyHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: service,
			Name:      "operations_latency_seconds",
			Help:      "latency of internal operations",
			Buckets:   []float64{0.01, 0.10, 0.50, 1.00, 2.00},
		},
		[]string{"op", "success"},
	)

	OpLatencySumm := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: service,
			Name:      "operations_latency_quantiles_seconds",
			Help:      "latency quantiles of internal operations",
			Objectives: map[float64]float64{
				0.1:  0.1,
				0.5:  0.05,
				0.95: 0.01,
				0.99: 0.001,
			},
		},
		[]string{"op", "success"},
	)

	registry.MustRegister(ReqCounter)
	registry.MustRegister(OpLatencySumm)
	registry.MustRegister(OpLatencyHist)

	return &Metrics{
		Registry:      registry,
		ReqCounter:    ReqCounter,
		OpLatencyHist: OpLatencyHist,
		OpLatencySumm: OpLatencySumm,
	}
}

// Handler returns http handler for metrics endpoint
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{})
}
