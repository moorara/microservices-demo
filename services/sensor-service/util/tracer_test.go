package util

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/moorara/microservices-demo/services/sensor-service/config"
	"github.com/stretchr/testify/assert"
)

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
		logger log.Logger
	}{
		{
			"WithSpanLogging",
			config.Config{
				ServiceName:            "go-service",
				JaegerAgentHost:        "localhost",
				JaegerAgentPort:        6831,
				JaegerReporterLogSpans: true,
			},
			log.NewNopLogger(),
		},
		{
			"WithoutSpanLogging",
			config.Config{
				ServiceName:            "go-service",
				JaegerAgentHost:        "localhost",
				JaegerAgentPort:        6831,
				JaegerReporterLogSpans: false,
			},
			log.NewNopLogger(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer, closer := NewTracer(tc.config, tc.logger)
			defer closer.Close()

			assert.NotNil(t, tracer)
			assert.NotNil(t, closer)
		})
	}
}
