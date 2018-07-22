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
	Registry *prometheus.Registry
}

// NewMetrics creates a Metrics instance
func NewMetrics(service string) *Metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), service))

	return &Metrics{
		service:  service,
		Registry: registry,
	}
}

// NewCounter creates and registers a new counter metric
func (m *Metrics) NewCounter(prefixName bool, name, help string, labels []string) *prometheus.CounterVec {
	opts := prometheus.CounterOpts{
		Name: name,
		Help: help,
	}

	if prefixName == true {
		opts.Namespace = m.service
	}

	counter := prometheus.NewCounterVec(opts, labels)
	m.Registry.MustRegister(counter)
	return counter
}

// NewGauge creates and registers a new gauge metric
func (m *Metrics) NewGauge(prefixName bool, name, help string, labels []string) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}

	if prefixName == true {
		opts.Namespace = m.service
	}

	gauge := prometheus.NewGaugeVec(opts, labels)
	m.Registry.MustRegister(gauge)
	return gauge
}

// NewHistogram creates and registers a new histogram metric
func (m *Metrics) NewHistogram(prefixName bool, name, help string, buckets []float64, labels []string) *prometheus.HistogramVec {
	opts := prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	}

	if prefixName == true {
		opts.Namespace = m.service
	}

	histogram := prometheus.NewHistogramVec(opts, labels)
	m.Registry.MustRegister(histogram)
	return histogram
}

// NewSummary creates and registers a new summary metric
func (m *Metrics) NewSummary(prefixName bool, name, help string, quantiles map[float64]float64, labels []string) *prometheus.SummaryVec {
	opts := prometheus.SummaryOpts{
		Name:       name,
		Help:       help,
		Objectives: quantiles,
	}

	if prefixName == true {
		opts.Namespace = m.service
	}

	summary := prometheus.NewSummaryVec(opts, labels)
	m.Registry.MustRegister(summary)
	return summary
}

// GetHandler returns http handler for metrics endpoint
func (m *Metrics) GetHandler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{})
}
