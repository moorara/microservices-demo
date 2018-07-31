package transport

import (
	"context"
	"testing"

	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/stretchr/testify/assert"
)

// mockSwitchService is a mock implementation of proto.SwitchServiceServer
type mockSwitchService struct {
	InstallSwitchCallCount int
	InstallSwitchInContext context.Context
	InstallSwitchInReq     *proto.InstallSwitchRequest
	InstallSwitchOutResp   *proto.Switch
	InstallSwitchOutError  error

	RemoveSwitchCallCount int
	RemoveSwitchInContext context.Context
	RemoveSwitchInReq     *proto.RemoveSwitchRequest
	RemoveSwitchOutResp   *proto.RemoveSwitchResponse
	RemoveSwitchOutError  error

	GetSwitchCallCount int
	GetSwitchInContext context.Context
	GetSwitchInReq     *proto.GetSwitchRequest
	GetSwitchOutResp   *proto.Switch
	GetSwitchOutError  error

	GetSwitchesCallCount int
	GetSwitchesInReq     *proto.GetSwitchesRequest
	GetSwitchesInStream  proto.SwitchService_GetSwitchesServer
	GetSwitchesOutError  error

	SetSwitchCallCount int
	SetSwitchInContext context.Context
	SetSwitchInReq     *proto.SetSwitchRequest
	SetSwitchOutResp   *proto.SetSwitchResponse
	SetSwitchOutError  error
}

func (m *mockSwitchService) InstallSwitch(ctx context.Context, req *proto.InstallSwitchRequest) (*proto.Switch, error) {
	m.InstallSwitchCallCount++
	m.InstallSwitchInContext = ctx
	m.InstallSwitchInReq = req
	return m.InstallSwitchOutResp, m.InstallSwitchOutError
}

func (m *mockSwitchService) RemoveSwitch(ctx context.Context, req *proto.RemoveSwitchRequest) (*proto.RemoveSwitchResponse, error) {
	m.RemoveSwitchCallCount++
	m.RemoveSwitchInContext = ctx
	m.RemoveSwitchInReq = req
	return m.RemoveSwitchOutResp, m.RemoveSwitchOutError
}

func (m *mockSwitchService) GetSwitch(ctx context.Context, req *proto.GetSwitchRequest) (*proto.Switch, error) {
	m.GetSwitchCallCount++
	m.GetSwitchInContext = ctx
	m.GetSwitchInReq = req
	return m.GetSwitchOutResp, m.GetSwitchOutError
}

func (m *mockSwitchService) GetSwitches(req *proto.GetSwitchesRequest, stream proto.SwitchService_GetSwitchesServer) error {
	m.GetSwitchesCallCount++
	m.GetSwitchesInReq = req
	m.GetSwitchesInStream = stream
	return m.GetSwitchesOutError
}

func (m *mockSwitchService) SetSwitch(ctx context.Context, req *proto.SetSwitchRequest) (*proto.SetSwitchResponse, error) {
	m.SetSwitchCallCount++
	m.SetSwitchInContext = ctx
	m.SetSwitchInReq = req
	return m.SetSwitchOutResp, m.SetSwitchOutError
}

func TestNewGRPCServer(t *testing.T) {
	tests := []struct {
		name          string
		caFile        string
		certFile      string
		keyFile       string
		switchService proto.SwitchServiceServer
		expectError   bool
	}{
		{
			name:          "Simple",
			switchService: &mockSwitchService{},
			expectError:   false,
		},
		{
			name:          "MTLS",
			caFile:        "../certs/ca.chain.cert",
			certFile:      "../certs/server.cert",
			keyFile:       "../certs/server.key",
			switchService: &mockSwitchService{},
			expectError:   false,
		},
		{
			name:          "NoCACert",
			caFile:        "../certs/ca.chain",
			certFile:      "../certs/server.cert",
			keyFile:       "../certs/server.key",
			switchService: &mockSwitchService{},
			expectError:   true,
		},
		{
			name:          "NoCert",
			caFile:        "../certs/ca.chain.cert",
			certFile:      "../certs/server",
			keyFile:       "../certs/server.key",
			switchService: &mockSwitchService{},
			expectError:   true,
		},
		{
			name:          "NoKey",
			caFile:        "../certs/ca.chain.cert",
			certFile:      "../certs/server.cert",
			keyFile:       "../certs/server",
			switchService: &mockSwitchService{},
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grpcServer, err := NewGRPCServer(tc.caFile, tc.certFile, tc.keyFile, tc.switchService)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, grpcServer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, grpcServer)
			}
		})
	}
}
