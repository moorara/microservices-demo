package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		service string
	}{
		{"service_name"},
		{"service-name"},
	}

	for _, tc := range tests {
		metrics := New(tc.service)
		handler := metrics.Handler()

		assert.NotNil(t, handler)
		assert.NotNil(t, metrics.Registry)
		assert.NotNil(t, metrics.OpLatencyHist)
		assert.NotNil(t, metrics.OpLatencySumm)
	}
}
