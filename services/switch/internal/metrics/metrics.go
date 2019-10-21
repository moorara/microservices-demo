package metrics

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics includes all metrics
type Metrics struct {
	Registry      *prometheus.Registry
	ReqCounter    *prometheus.CounterVec
	OpLatencyHist *prometheus.HistogramVec
	OpLatencySumm *prometheus.SummaryVec
}

// Mock creates a new metrics for testing purposes
func Mock() *Metrics {
	registry := prometheus.NewRegistry()

	ReqCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
		},
		[]string{"call", "success"},
	)

	OpLatencyHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "operations_latency_seconds",
		},
		[]string{"op", "success"},
	)

	OpLatencySumm := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "operations_latency_quantiles_seconds",
		},
		[]string{"op", "success"},
	)

	return &Metrics{
		Registry:      registry,
		ReqCounter:    ReqCounter,
		OpLatencyHist: OpLatencyHist,
		OpLatencySumm: OpLatencySumm,
	}
}

// New creates a new metrics
func New(service string) *Metrics {
	service = strings.Replace(service, "-", "_", -1)

	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
		Namespace: service,
	}))

	// ReqCounter is a counter tracking the total number of requests
	ReqCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: service,
			Name:      "grpc_requests_total",
			Help:      "total number of grpc requests",
		},
		[]string{"call", "success"},
	)

	// OpLatencyHist is a histogram tracking the response times of internal operations
	OpLatencyHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: service,
			Name:      "operations_latency_seconds",
			Help:      "latency of internal operations",
			Buckets:   []float64{0.01, 0.10, 0.50, 1.00},
		},
		[]string{"op", "success"},
	)

	// OpLatencySumm is a summary tracking the response times of internal operations
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
