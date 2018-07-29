package service

import (
	"context"

	arango "github.com/arangodb/go-driver"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"google.golang.org/grpc/metadata"
)

// MockServerStream is a mock implementation of grpc.ServerStream
type MockServerStream struct {
	SetHeaderCallCount int
	SetHeaderInMeta    metadata.MD
	SetHeaderOutError  error

	SendHeaderCallCount int
	SendHeaderInMeta    metadata.MD
	SendHeaderOutError  error

	SetTrailerCallCount int
	SetTrailerInMeta    metadata.MD

	ContextCallCount  int
	ContextOutContext context.Context

	SendMsgCallCount int
	SendMsgInMsg     interface{}
	SendMsgOutError  error

	RecvMsgCallCount int
	RecvMsgInMsg     interface{}
	RecvMsgOutError  error
}

func (m *MockServerStream) SetHeader(md metadata.MD) error {
	m.SetHeaderCallCount++
	m.SetHeaderInMeta = md
	return m.SetHeaderOutError
}

func (m *MockServerStream) SendHeader(md metadata.MD) error {
	m.SendHeaderCallCount++
	m.SendHeaderInMeta = md
	return m.SendHeaderOutError
}

func (m *MockServerStream) SetTrailer(md metadata.MD) {
	m.SetTrailerCallCount++
	m.SetTrailerInMeta = md
}

func (m *MockServerStream) Context() context.Context {
	m.ContextCallCount++
	return m.ContextOutContext
}

func (m *MockServerStream) SendMsg(msg interface{}) error {
	m.SendMsgCallCount++
	m.SendMsgInMsg = msg
	return m.SendMsgOutError
}

func (m *MockServerStream) RecvMsg(msg interface{}) error {
	m.RecvMsgCallCount++
	m.RecvMsgInMsg = msg
	return m.RecvMsgOutError
}

// MockArangoService is a mock implementation of service.ArangoService
type MockArangoService struct {
	QueryCallCount int
	QueryInContext context.Context
	QueryInQuery   string
	QueryInVars    map[string]interface{}
	QueryOutCursor arango.Cursor
	QueryOutError  error

	CreateDocumentCallCount int
	CreateDocumentInContext context.Context
	CreateDocumentInDoc     interface{}
	CreateDocumentOutMeta   arango.DocumentMeta
	CreateDocumentOutError  error

	ReadDocumentCallCount int
	ReadDocumentInContext context.Context
	ReadDocumentInKey     string
	ReadDocumentInDoc     interface{}
	ReadDocumentOutMeta   arango.DocumentMeta
	ReadDocumentOutError  error

	UpdateDocumentCallCount int
	UpdateDocumentInContext context.Context
	UpdateDocumentInKey     string
	UpdateDocumentInDoc     interface{}
	UpdateDocumentOutMeta   arango.DocumentMeta
	UpdateDocumentOutError  error

	RemoveDocumentCallCount int
	RemoveDocumentInContext context.Context
	RemoveDocumentInKey     string
	RemoveDocumentOutMeta   arango.DocumentMeta
	RemoveDocumentOutError  error
}

func (m *MockArangoService) Query(ctx context.Context, query string, vars map[string]interface{}) (arango.Cursor, error) {
	m.QueryCallCount++
	m.QueryInContext = ctx
	m.QueryInQuery = query
	m.QueryInVars = vars
	return m.QueryOutCursor, m.QueryOutError
}

func (m *MockArangoService) CreateDocument(ctx context.Context, doc interface{}) (arango.DocumentMeta, error) {
	m.CreateDocumentCallCount++
	m.CreateDocumentInContext = ctx
	m.CreateDocumentInDoc = doc
	return m.CreateDocumentOutMeta, m.CreateDocumentOutError
}

func (m *MockArangoService) ReadDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.ReadDocumentCallCount++
	m.ReadDocumentInContext = ctx
	m.ReadDocumentInKey = key
	m.ReadDocumentInDoc = doc
	return m.ReadDocumentOutMeta, m.ReadDocumentOutError
}

func (m *MockArangoService) UpdateDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.UpdateDocumentCallCount++
	m.UpdateDocumentInContext = ctx
	m.UpdateDocumentInKey = key
	m.UpdateDocumentInDoc = doc
	return m.UpdateDocumentOutMeta, m.UpdateDocumentOutError
}

func (m *MockArangoService) RemoveDocument(ctx context.Context, key string) (arango.DocumentMeta, error) {
	m.RemoveDocumentCallCount++
	m.RemoveDocumentInContext = ctx
	m.RemoveDocumentInKey = key
	return m.RemoveDocumentOutMeta, m.RemoveDocumentOutError
}

// MockSwitchService is a mock implementation of proto.SwitchServiceServer
type MockSwitchService struct {
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

func (m *MockSwitchService) InstallSwitch(ctx context.Context, req *proto.InstallSwitchRequest) (*proto.Switch, error) {
	m.InstallSwitchCallCount++
	m.InstallSwitchInContext = ctx
	m.InstallSwitchInReq = req
	return m.InstallSwitchOutResp, m.InstallSwitchOutError
}

func (m *MockSwitchService) RemoveSwitch(ctx context.Context, req *proto.RemoveSwitchRequest) (*proto.RemoveSwitchResponse, error) {
	m.RemoveSwitchCallCount++
	m.RemoveSwitchInContext = ctx
	m.RemoveSwitchInReq = req
	return m.RemoveSwitchOutResp, m.RemoveSwitchOutError
}

func (m *MockSwitchService) GetSwitch(ctx context.Context, req *proto.GetSwitchRequest) (*proto.Switch, error) {
	m.GetSwitchCallCount++
	m.GetSwitchInContext = ctx
	m.GetSwitchInReq = req
	return m.GetSwitchOutResp, m.GetSwitchOutError
}

func (m *MockSwitchService) GetSwitches(req *proto.GetSwitchesRequest, stream proto.SwitchService_GetSwitchesServer) error {
	m.GetSwitchesCallCount++
	m.GetSwitchesInReq = req
	m.GetSwitchesInStream = stream
	return m.GetSwitchesOutError
}

func (m *MockSwitchService) SetSwitch(ctx context.Context, req *proto.SetSwitchRequest) (*proto.SetSwitchResponse, error) {
	m.SetSwitchCallCount++
	m.SetSwitchInContext = ctx
	m.SetSwitchInReq = req
	return m.SetSwitchOutResp, m.SetSwitchOutError
}
