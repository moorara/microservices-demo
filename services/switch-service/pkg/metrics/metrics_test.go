package metrics

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getMetricString(t *testing.T, metrics *Metrics, namespace, service string) (output string) {
	metricFamilies, err := metrics.Registry.Gather()
	assert.NoError(t, err)

	for _, mf := range metricFamilies {
		if strings.HasSuffix(*mf.Name, service) {
			for _, m := range mf.Metric {
				output += m.String() + "\n"
			}
		}
	}

	return output
}

func TestNewVoidMetrics(t *testing.T) {
	metrics := NewVoidMetrics()
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.Registry)
}

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		service string
	}{
		{"service_name"},
		{"service-name"},
	}

	for _, tc := range tests {
		metrics := NewMetrics(tc.service)
		assert.NotNil(t, metrics)
	}
}

func TestNewCounter(t *testing.T) {
	tests := []struct {
		name          string
		service       string
		metric, help  string
		labels        []string
		labelValues   []string
		addValue      float64
		expectedRegex []string
	}{
		{
			"ErrorCounter",
			"service_name",
			"service_name_errors_total", "total number of errors",
			[]string{"resource"},
			[]string{"object"},
			5,
			[]string{
				`label:<name:"resource" value:"object" >`,
			},
		},
		{
			"RequestCounter",
			"service_name",
			"service_name_requests_total", "total number of requests",
			[]string{"resource"},
			[]string{"object"},
			10,
			[]string{
				`label:<name:"resource" value:"object" >`,
				`counter:<value:11 > `,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metrics := NewMetrics(tc.service)
			counter := metrics.NewCounter(tc.metric, tc.help, tc.labels)
			counter.WithLabelValues(tc.labelValues...).Inc()
			counter.WithLabelValues(tc.labelValues...).Add(tc.addValue)
			output := getMetricString(t, metrics, tc.service, tc.metric)

			for _, rx := range tc.expectedRegex {
				assert.Regexp(t, rx, output)
			}
		})
	}
}

func TestNewGauge(t *testing.T) {
	tests := []struct {
		name               string
		service            string
		metric, help       string
		labels             []string
		labelValues        []string
		addValue, subValue float64
		expectedRegex      []string
	}{
		{
			"ObjectGauge",
			"service_name",
			"service_name_objects", "current number of objects",
			[]string{"active"},
			[]string{"true"},
			8, 4,
			[]string{
				`label:<name:"active" value:"true" >`,
				`gauge:<value:4 > `,
			},
		},
		{
			"ConnectionGauge",
			"service_name",
			"service_name_connections", "active number of connections",
			[]string{"type"},
			[]string{"tcp"},
			10, 5,
			[]string{
				`label:<name:"type" value:"tcp" >`,
				`gauge:<value:5 >`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metrics := NewMetrics(tc.service)
			gauge := metrics.NewGauge(tc.metric, tc.help, tc.labels)
			gauge.WithLabelValues(tc.labelValues...).Inc()
			gauge.WithLabelValues(tc.labelValues...).Dec()
			gauge.WithLabelValues(tc.labelValues...).Add(tc.addValue)
			gauge.WithLabelValues(tc.labelValues...).Sub(tc.subValue)
			output := getMetricString(t, metrics, tc.service, tc.metric)

			for _, rx := range tc.expectedRegex {
				assert.Regexp(t, rx, output)
			}
		})
	}
}

func TestNewHistogram(t *testing.T) {
	tests := []struct {
		name          string
		service       string
		metric, help  string
		buckets       []float64
		labels        []string
		labelValues   []string
		value         float64
		expectedRegex []string
	}{
		{
			"DurationHistogram",
			"service_name",
			"service_name_op_duration_seconds", "operation durations in seconds",
			[]float64{0.01, 0.1, 1},
			[]string{"op", "success"},
			[]string{"creation", "false"},
			0.1234,
			[]string{
				`label:<name:"op" value:"creation" >`,
				`label:<name:"success" value:"false" >`,
				`histogram:<sample_count:1 sample_sum:0.1234`,
				`bucket:<cumulative_count:0 upper_bound:0.01 >`,
				`bucket:<cumulative_count:0 upper_bound:0.1 >`,
				`bucket:<cumulative_count:1 upper_bound:1 >`,
			},
		},
		{
			"ThroughputHistogram",
			"service_name",
			"service_name_throughput_bytes_per_second", "operation throughput in bytes per second",
			[]float64{0.01, 0.1, 0.5, 1, 5, 10},
			[]string{"op", "success"},
			[]string{"deletion", "true"},
			1.666,
			[]string{
				`label:<name:"op" value:"deletion" >`,
				`label:<name:"success" value:"true" >`,
				`histogram:<sample_count:1 sample_sum:1.666`,
				`bucket:<cumulative_count:0 upper_bound:0.01 >`,
				`bucket:<cumulative_count:0 upper_bound:0.1 >`,
				`bucket:<cumulative_count:0 upper_bound:0.5 >`,
				`bucket:<cumulative_count:0 upper_bound:1 >`,
				`bucket:<cumulative_count:1 upper_bound:5 >`,
				`bucket:<cumulative_count:1 upper_bound:10 >`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metrics := NewMetrics(tc.service)
			histogram := metrics.NewHistogram(tc.metric, tc.help, tc.buckets, tc.labels)
			histogram.WithLabelValues(tc.labelValues...).Observe(tc.value)
			output := getMetricString(t, metrics, tc.service, tc.metric)

			for _, rx := range tc.expectedRegex {
				assert.Regexp(t, rx, output)
			}
		})
	}
}

func TestNewSummary(t *testing.T) {
	tests := []struct {
		name          string
		service       string
		metric, help  string
		quantiles     map[float64]float64
		labels        []string
		labelValues   []string
		value         float64
		expectedRegex []string
	}{
		{
			"DurationSummary",
			"service_name",
			"service_name_op_duration_seconds", "operation durations in seconds",
			map[float64]float64{0.5: 0.05, 0.9: 0.01},
			[]string{"op", "success"},
			[]string{"creation", "false"},
			0.1234,
			[]string{
				`label:<name:"op" value:"creation" >`,
				`label:<name:"success" value:"false" >`,
				`summary:<sample_count:1 sample_sum:0.1234`,
				`quantile:<quantile:0.5 value:0.1234 >`,
				`quantile:<quantile:0.9 value:0.1234 >`,
			},
		},
		{
			"ThroughputSummary",
			"service_name",
			"service_name_throughput_bytes_per_second", "operation throughput in bytes per second",
			map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			[]string{"op", "success"},
			[]string{"deletion", "true"},
			1.666,
			[]string{
				`label:<name:"op" value:"deletion" >`,
				`label:<name:"success" value:"true" >`,
				`summary:<sample_count:1 sample_sum:1.666`,
				`quantile:<quantile:0.5 value:1.666 >`,
				`quantile:<quantile:0.9 value:1.666 >`,
				`quantile:<quantile:0.99 value:1.666 >`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metrics := NewMetrics(tc.service)
			summmary := metrics.NewSummary(tc.metric, tc.help, tc.quantiles, tc.labels)
			summmary.WithLabelValues(tc.labelValues...).Observe(tc.value)
			output := getMetricString(t, metrics, tc.service, tc.metric)

			for _, rx := range tc.expectedRegex {
				assert.Regexp(t, rx, output)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name          string
		service       string
		prepare       func(m *Metrics)
		expectedRegex []string
	}{
		{
			"Histogram",
			"service_name",
			func(m *Metrics) {
				histogram := m.NewHistogram("service_name_op_duration_seconds", "operation durations", []float64{0.01, 0.1, 1}, []string{"op", "success"})
				histogram.WithLabelValues("creation", "true").Observe(0.27)
			},
			[]string{
				`# HELP go_[A-Za-z_]+ Number of`,
				`# TYPE go_[A-Za-z_]+ counter`,
				`# TYPE go_[A-Za-z_]+ gauge`,
				`go_[A-Za-z_]+ [-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`,
				`# TYPE service_name_op_duration_seconds histogram`,
				`service_name_op_duration_seconds_bucket{op="creation",success="true",le="0.01"} 0`,
				`service_name_op_duration_seconds_bucket{op="creation",success="true",le="0.1"} 0`,
				`service_name_op_duration_seconds_bucket{op="creation",success="true",le="1"} 1`,
				`service_name_op_duration_seconds_bucket{op="creation",success="true",le="\+Inf"} 1`,
				`service_name_op_duration_seconds_sum{op="creation",success="true"} 0.27`,
				`service_name_op_duration_seconds_count{op="creation",success="true"} 1`,
			},
		},
		{
			"Summary",
			"service_name",
			func(m *Metrics) {
				summary := m.NewSummary("service_name_op_duration_seconds", "operation durations", map[float64]float64{0.5: 0.05, 0.9: 0.01}, []string{"op", "success"})
				summary.WithLabelValues("creation", "true").Observe(0.27)
			},
			[]string{
				`# HELP go_[A-Za-z_]+ Number of`,
				`# TYPE go_[A-Za-z_]+ counter`,
				`# TYPE go_[A-Za-z_]+ gauge`,
				`go_[A-Za-z_]+ [-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`,
				`# TYPE service_name_op_duration_seconds summary`,
				`service_name_op_duration_seconds{op="creation",success="true",quantile="0.5"} 0.27`,
				`service_name_op_duration_seconds{op="creation",success="true",quantile="0.9"} 0.27`,
				`service_name_op_duration_seconds_sum{op="creation",success="true"} 0.27`,
				`service_name_op_duration_seconds_count{op="creation",success="true"} 1`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metrics := NewMetrics(tc.service)
			tc.prepare(metrics)
			handler := metrics.Handler()

			ts := httptest.NewServer(handler)
			defer ts.Close()
			res, err := http.Get(ts.URL)
			assert.NoError(t, err)
			content, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			output := string(content)

			for _, rx := range tc.expectedRegex {
				assert.Regexp(t, rx, output)
			}
		})
	}
}
