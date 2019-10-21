package util

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moorara/microservices-demo/services/sensor/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name       string
		config     config.Config
		logger     log.Logger
		registerer prometheus.Registerer
	}{
		{
			"WithSpanLogging",
			config.Config{
				ServiceName:     "go-service",
				JaegerAgentAddr: "localhost:6831",
				JaegerLogSpans:  true,
			},
			log.NewNopLogger(),
			prometheus.NewRegistry(),
		},
		{
			"WithoutSpanLogging",
			config.Config{
				ServiceName:     "go-service",
				JaegerAgentAddr: "localhost:6831",
				JaegerLogSpans:  false,
			},
			log.NewNopLogger(),
			prometheus.NewRegistry(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer, closer := NewTracer(tc.config, tc.logger, tc.registerer)
			defer closer.Close()

			assert.NotNil(t, tracer)
			assert.NotNil(t, closer)
		})
	}
}
