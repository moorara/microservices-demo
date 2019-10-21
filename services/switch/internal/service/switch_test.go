package service

import (
	"context"
	"errors"
	"testing"

	"github.com/moorara/microservices-demo/services/switch/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch/internal/proto"
	"github.com/moorara/microservices-demo/services/switch/pkg/log"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"

	arango "github.com/arangodb/go-driver"
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
		name           string
		arango         ArangoService
		ctx            context.Context
		req            *proto.InstallSwitchRequest
		expectedError  error
		expectedSwitch *proto.Switch
	}{
		{
			"Fail",
			&mockArangoService{
				CreateDocumentOutError: errors.New("database error"),
			},
			context.Background(),
			&proto.InstallSwitchRequest{
				SiteId: "1111-1111",
				Name:   "Light",
				State:  "OFF",
				States: []string{"ON", "OFF"},
			},
			errors.New("database error"),
			nil,
		},
		{
			"Success",
			&mockArangoService{
				CreateDocumentOutMeta: arango.DocumentMeta{
					Key: "aaaa-aaaa",
				},
			},
			context.Background(),
			&proto.InstallSwitchRequest{
				SiteId: "1111-1111",
				Name:   "Light",
				State:  "OFF",
				States: []string{"ON", "OFF"},
			},
			nil,
			&proto.Switch{
				Id:     "aaaa-aaaa",
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

			sw, err := service.InstallSwitch(tc.ctx, tc.req)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedSwitch, sw)
		})
	}
}

func TestRemoveSwitch(t *testing.T) {
	tests := []struct {
		name             string
		arango           ArangoService
		ctx              context.Context
		req              *proto.RemoveSwitchRequest
		expectedError    error
		expectedResponse *proto.RemoveSwitchResponse
	}{
		{
			"Fail",
			&mockArangoService{
				RemoveDocumentOutError: errors.New("database error"),
			},
			context.Background(),
			&proto.RemoveSwitchRequest{
				Id: "aaaa-aaaa",
			},
			errors.New("database error"),
			nil,
		},
		{
			"Success",
			&mockArangoService{
				RemoveDocumentOutMeta: arango.DocumentMeta{},
			},
			context.Background(),
			&proto.RemoveSwitchRequest{
				Id: "aaaa-aaaa",
			},
			nil,
			&proto.RemoveSwitchResponse{},
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

			resp, err := service.RemoveSwitch(tc.ctx, tc.req)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

func TestGetSwitch(t *testing.T) {
	tests := []struct {
		name           string
		arango         ArangoService
		ctx            context.Context
		req            *proto.GetSwitchRequest
		expectedError  error
		expectedSwitch *proto.Switch
	}{
		{
			"Fail",
			&mockArangoService{
				ReadDocumentOutError: errors.New("database error"),
			},
			context.Background(),
			&proto.GetSwitchRequest{
				Id: "aaaa-aaaa",
			},
			errors.New("database error"),
			nil,
		},
		{
			"Success",
			&mockArangoService{
				ReadDocumentOutMeta: arango.DocumentMeta{
					Key: "aaaa-aaaa",
				},
			},
			context.Background(),
			&proto.GetSwitchRequest{
				Id: "aaaa-aaaa",
			},
			nil,
			&proto.Switch{},
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

			sw, err := service.GetSwitch(tc.ctx, tc.req)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedSwitch, sw)
		})
	}
}

func TestGetSwitches(t *testing.T) {
	tests := []struct {
		name          string
		arango        ArangoService
		req           *proto.GetSwitchesRequest
		stream        *mockGetSwitchesServer
		expectedError error
	}{
		{
			"QueryError",
			&mockArangoService{
				QueryOutError: errors.New("database error"),
			},
			&proto.GetSwitchesRequest{
				SiteId: "1111-1111",
			},
			&mockGetSwitchesServer{
				ServerStream: &mockServerStream{
					ContextOutContext: context.Background(),
				},
			},
			errors.New("database error"),
		},
		{
			"ReadDocumentError",
			&mockArangoService{
				QueryOutCursor: &mockArangoCursor{
					Closer:               &mockCloser{},
					HasMoreOutResults:    []bool{true},
					ReadDocumentOutError: errors.New("cursor error"),
				},
			},
			&proto.GetSwitchesRequest{
				SiteId: "1111-1111",
			},
			&mockGetSwitchesServer{
				ServerStream: &mockServerStream{
					ContextOutContext: context.Background(),
				},
			},
			errors.New("cursor error"),
		},
		{
			"SendError",
			&mockArangoService{
				QueryOutCursor: &mockArangoCursor{
					Closer:            &mockCloser{},
					HasMoreOutResults: []bool{true},
					ReadDocumentOutMeta: arango.DocumentMeta{
						Key: "aaaa-aaaa",
					},
				},
			},
			&proto.GetSwitchesRequest{
				SiteId: "1111-1111",
			},
			&mockGetSwitchesServer{
				ServerStream: &mockServerStream{
					ContextOutContext: context.Background(),
				},
				SendOutError: errors.New("stream error"),
			},
			errors.New("stream error"),
		},
		{
			"Success",
			&mockArangoService{
				QueryOutCursor: &mockArangoCursor{
					Closer:            &mockCloser{},
					HasMoreOutResults: []bool{true, false},
					ReadDocumentOutMeta: arango.DocumentMeta{
						Key: "aaaa-aaaa",
					},
				},
			},
			&proto.GetSwitchesRequest{
				SiteId: "1111-1111",
			},
			&mockGetSwitchesServer{
				ServerStream: &mockServerStream{
					ContextOutContext: context.Background(),
				},
				SendOutError: nil,
			},
			nil,
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

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestSetSwitch(t *testing.T) {
	tests := []struct {
		name             string
		arango           ArangoService
		ctx              context.Context
		req              *proto.SetSwitchRequest
		expectedError    error
		expectedResponse *proto.SetSwitchResponse
	}{
		{
			"Fail",
			&mockArangoService{
				UpdateDocumentOutError: errors.New("database error"),
			},
			context.Background(),
			&proto.SetSwitchRequest{
				Id:    "aaaa-aaaa",
				State: "ON",
			},
			errors.New("database error"),
			nil,
		},
		{
			"Success",
			&mockArangoService{
				UpdateDocumentOutMeta: arango.DocumentMeta{},
			},
			context.Background(),
			&proto.SetSwitchRequest{
				Id:    "aaaa-aaaa",
				State: "ON",
			},
			nil,
			&proto.SetSwitchResponse{},
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

			resp, err := service.SetSwitch(tc.ctx, tc.req)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
