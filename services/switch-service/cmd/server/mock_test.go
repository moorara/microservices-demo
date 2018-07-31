package server

import (
	"context"
	"net"
	"net/http"

	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
)

type mockHTTPServer struct {
	ListenAndServeCallCount int
	ListenAndServeOutError  error

	ServeCallCount  int
	ServeInListener net.Listener
	ServeOutError   error

	CloseCallCount int
	CloseOutError  error

	ShutdownCallCount int
	ShutdownInContext context.Context
	ShutdownOutError  error
}

func (m *mockHTTPServer) ListenAndServe() error {
	m.ListenAndServeCallCount++
	return m.ListenAndServeOutError
}

func (m *mockHTTPServer) Serve(listener net.Listener) error {
	m.ServeCallCount++
	m.ServeInListener = listener
	return m.ServeOutError
}

func (m *mockHTTPServer) Close() error {
	m.CloseCallCount++
	return m.CloseOutError
}

func (m *mockHTTPServer) Shutdown(ctx context.Context) error {
	m.ShutdownCallCount++
	m.ShutdownInContext = ctx
	return m.ShutdownOutError
}

type mockGRPCServer struct {
	ServeCallCount  int
	ServeInListener net.Listener
	ServeOutError   error

	ServeHTTPCallCount int
	ServeHTTPInResp    http.ResponseWriter
	ServeHTTPInReq     *http.Request

	StopCallCount int

	GracefulStopCallCount int
}

func (m *mockGRPCServer) Serve(listener net.Listener) error {
	if listener != nil {
		listener.Close()
	}

	m.ServeCallCount++
	m.ServeInListener = listener
	return m.ServeOutError
}

func (m *mockGRPCServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	m.ServeHTTPCallCount++
	m.ServeHTTPInResp = resp
	m.ServeHTTPInReq = req
}

func (m *mockGRPCServer) Stop() {
	m.StopCallCount++
}

func (m *mockGRPCServer) GracefulStop() {
	m.GracefulStopCallCount++
}

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
