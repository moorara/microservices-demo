package util

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics represents metrics utility
type Metrics struct {
	service  string
	registry *prometheus.Registry
}

// NewMetrics creates a Metrics instance
func NewMetrics(service string) *Metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), service))

	return &Metrics{
		service:  service,
		registry: registry,
	}
}

// NewCounter creates and registers a new counter metric
func (m *Metrics) NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.service,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	m.registry.MustRegister(counter)
	return counter
}

// NewGauge creates and registers a new gauge metric
func (m *Metrics) NewGauge(name, help string, labels []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.service,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	m.registry.MustRegister(gauge)
	return gauge
}

// NewHistogram creates and registers a new histogram metric
func (m *Metrics) NewHistogram(name, help string, buckets []float64, labels []string) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: m.service,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labels,
	)

	m.registry.MustRegister(histogram)
	return histogram
}

// NewSummary creates and registers a new summary metric
func (m *Metrics) NewSummary(name, help string, quantiles map[float64]float64, labels []string) *prometheus.SummaryVec {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  m.service,
			Name:       name,
			Help:       help,
			Objectives: quantiles,
		},
		labels,
	)

	m.registry.MustRegister(summary)
	return summary
}

// GetHandler returns http handler for metrics endpoint
func (m *Metrics) GetHandler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
