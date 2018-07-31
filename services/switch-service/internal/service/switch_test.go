package service

import (
	"context"
	"testing"

	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestNewSwitchService(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
	}{
		{
			"Default",
			&mockArangoService{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			service := NewSwitchService(tc.arango, logger, metrics, tracer)

			assert.NotNil(t, service)
		})
	}
}

func TestInstallSwitch(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
		req    *proto.InstallSwitchRequest
	}{
		{
			"Simple",
			&mockArangoService{},
			&proto.InstallSwitchRequest{
				SiteId: "1111-1111",
				Name:   "Light",
				State:  "OFF",
				States: []string{"ON", "OFF"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			service := &SwitchService{
				arango:  tc.arango,
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			response, err := service.InstallSwitch(context.Background(), tc.req)

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestRemoveSwitch(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
		req    *proto.RemoveSwitchRequest
	}{
		{
			"Simple",
			&mockArangoService{},
			&proto.RemoveSwitchRequest{
				Id: "aaaa-aaaa",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			service := &SwitchService{
				arango:  tc.arango,
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			response, err := service.RemoveSwitch(context.Background(), tc.req)

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestGetSwitch(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
		req    *proto.GetSwitchRequest
	}{
		{
			"Simple",
			&mockArangoService{},
			&proto.GetSwitchRequest{
				Id: "aaaa-aaaa",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			service := &SwitchService{
				arango:  tc.arango,
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			response, err := service.GetSwitch(context.Background(), tc.req)

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}

func TestGetSwitches(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
		req    *proto.GetSwitchesRequest
		stream *mockGetSwitchesServer
	}{
		{
			"Simple",
			&mockArangoService{},
			&proto.GetSwitchesRequest{
				SiteId: "1111-1111",
			},
			&mockGetSwitchesServer{
				SendOutError: nil,
				ServerStream: &mockServerStream{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			service := &SwitchService{
				arango:  tc.arango,
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			err := service.GetSwitches(tc.req, tc.stream)

			assert.NoError(t, err)
			assert.True(t, tc.stream.SendCalled)
		})
	}
}

func TestSetSwitch(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
		req    *proto.SetSwitchRequest
	}{
		{
			"Simple",
			&mockArangoService{},
			&proto.SetSwitchRequest{
				Id:    "aaaa-aaaa",
				State: "ON",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewVoidLogger()
			metrics := metrics.Mock()
			tracer := mocktracer.New()
			service := &SwitchService{
				arango:  tc.arango,
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			response, err := service.SetSwitch(context.Background(), tc.req)

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	}
}
