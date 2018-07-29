package service

import (
	"context"
	"testing"

	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// mockGetSwitchesServer mocks proto.SwitchService_GetSwitchesServer
type mockGetSwitchesServer struct {
	grpc.ServerStream
	SendCallCount int
	SendInSwitch  *proto.Switch
	SendOutError  error
}

func (m *mockGetSwitchesServer) Send(sw *proto.Switch) error {
	m.SendCallCount++
	m.SendInSwitch = sw
	return m.SendOutError
}

func TestNewSwitchService(t *testing.T) {
	tests := []struct {
		name   string
		arango ArangoService
	}{
		{
			"Default",
			&MockArangoService{},
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
			&MockArangoService{},
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
			&MockArangoService{},
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
			&MockArangoService{},
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
		name                  string
		arango                ArangoService
		req                   *proto.GetSwitchesRequest
		stream                *mockGetSwitchesServer
		expectedSendCallCount int
	}{
		{
			"Simple",
			&MockArangoService{},
			&proto.GetSwitchesRequest{
				SiteId: "1111-1111",
			},
			&mockGetSwitchesServer{
				SendOutError: nil,
				ServerStream: &MockServerStream{},
			},
			1,
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
			assert.Equal(t, tc.expectedSendCallCount, tc.stream.SendCallCount)
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
			&MockArangoService{},
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
