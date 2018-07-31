package server

import (
	"context"
	"net"
	"net/http"

	arango "github.com/arangodb/go-driver"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
)

// mockHTTPServer is a mock implementation of transport.HTTPServer
type mockHTTPServer struct {
	ListenAndServeCalled   bool
	ListenAndServeOutError error

	ServeCalled     bool
	ServeInListener net.Listener
	ServeOutError   error

	CloseCalled   bool
	CloseOutError error

	ShutdownCalled    bool
	ShutdownInContext context.Context
	ShutdownOutError  error
}

func (m *mockHTTPServer) ListenAndServe() error {
	m.ListenAndServeCalled = true
	return m.ListenAndServeOutError
}

func (m *mockHTTPServer) Serve(listener net.Listener) error {
	m.ServeCalled = true
	m.ServeInListener = listener
	return m.ServeOutError
}

func (m *mockHTTPServer) Close() error {
	m.CloseCalled = true
	return m.CloseOutError
}

func (m *mockHTTPServer) Shutdown(ctx context.Context) error {
	m.ShutdownCalled = true
	m.ShutdownInContext = ctx
	return m.ShutdownOutError
}

// mockGRPCServer is a mock implementation of transport.GRPCServer
type mockGRPCServer struct {
	ServeCalled     bool
	ServeInListener net.Listener
	ServeOutError   error

	ServeHTTPCalled bool
	ServeHTTPInResp http.ResponseWriter
	ServeHTTPInReq  *http.Request

	StopCalled bool

	GracefulStopCalled bool
}

func (m *mockGRPCServer) Serve(listener net.Listener) error {
	if listener != nil {
		listener.Close()
	}

	m.ServeCalled = true
	m.ServeInListener = listener
	return m.ServeOutError
}

func (m *mockGRPCServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	m.ServeHTTPCalled = true
	m.ServeHTTPInResp = resp
	m.ServeHTTPInReq = req
}

func (m *mockGRPCServer) Stop() {
	m.StopCalled = true
}

func (m *mockGRPCServer) GracefulStop() {
	m.GracefulStopCalled = true
}

// mockArangoService is a mock implementation of service.ArangoService
type mockArangoService struct {
	ConnectCalled       bool
	ConnectInContext    context.Context
	ConnectInDatabase   string
	ConnectInCollection string
	ConnectOutError     error

	QueryCalled    bool
	QueryInContext context.Context
	QueryInQuery   string
	QueryInVars    map[string]interface{}
	QueryOutCursor arango.Cursor
	QueryOutError  error

	CreateDocumentCalled          bool
	CreateDocumentInContext       context.Context
	CreateDocumentInDoc           interface{}
	CreateDocumentOutDocumentMeta arango.DocumentMeta
	CreateDocumentOutError        error

	ReadDocumentCalled          bool
	ReadDocumentInContext       context.Context
	ReadDocumentInKey           string
	ReadDocumentInDoc           interface{}
	ReadDocumentOutDocumentMeta arango.DocumentMeta
	ReadDocumentOutError        error

	UpdateDocumentCalled          bool
	UpdateDocumentInContext       context.Context
	UpdateDocumentInKey           string
	UpdateDocumentInDoc           interface{}
	UpdateDocumentOutDocumentMeta arango.DocumentMeta
	UpdateDocumentOutError        error

	RemoveDocumentCalled          bool
	RemoveDocumentInContext       context.Context
	RemoveDocumentInKey           string
	RemoveDocumentOutDocumentMeta arango.DocumentMeta
	RemoveDocumentOutError        error
}

func (m *mockArangoService) Connect(ctx context.Context, database, collection string) error {
	m.ConnectCalled = true
	m.ConnectInContext = ctx
	m.ConnectInDatabase = database
	m.ConnectInCollection = collection
	return m.ConnectOutError
}

func (m *mockArangoService) Query(ctx context.Context, query string, vars map[string]interface{}) (arango.Cursor, error) {
	m.QueryCalled = true
	m.QueryInContext = ctx
	m.QueryInQuery = query
	m.QueryInVars = vars
	return m.QueryOutCursor, m.QueryOutError
}

func (m *mockArangoService) CreateDocument(ctx context.Context, doc interface{}) (arango.DocumentMeta, error) {
	m.CreateDocumentCalled = true
	m.CreateDocumentInContext = ctx
	m.CreateDocumentInDoc = doc
	return m.CreateDocumentOutDocumentMeta, m.CreateDocumentOutError
}

func (m *mockArangoService) ReadDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.ReadDocumentCalled = true
	m.ReadDocumentInContext = ctx
	m.ReadDocumentInKey = key
	m.ReadDocumentInDoc = doc
	return m.ReadDocumentOutDocumentMeta, m.ReadDocumentOutError
}

func (m *mockArangoService) UpdateDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.UpdateDocumentCalled = true
	m.UpdateDocumentInContext = ctx
	m.UpdateDocumentInKey = key
	m.UpdateDocumentInDoc = doc
	return m.UpdateDocumentOutDocumentMeta, m.UpdateDocumentOutError
}

func (m *mockArangoService) RemoveDocument(ctx context.Context, key string) (arango.DocumentMeta, error) {
	m.RemoveDocumentCalled = true
	m.RemoveDocumentInContext = ctx
	m.RemoveDocumentInKey = key
	return m.RemoveDocumentOutDocumentMeta, m.RemoveDocumentOutError
}

// mockSwitchService is a mock implementation of proto.SwitchServiceServer
type mockSwitchService struct {
	InstallSwitchCalled    bool
	InstallSwitchInContext context.Context
	InstallSwitchInReq     *proto.InstallSwitchRequest
	InstallSwitchOutResp   *proto.Switch
	InstallSwitchOutError  error

	RemoveSwitchCalled    bool
	RemoveSwitchInContext context.Context
	RemoveSwitchInReq     *proto.RemoveSwitchRequest
	RemoveSwitchOutResp   *proto.RemoveSwitchResponse
	RemoveSwitchOutError  error

	GetSwitchCalled    bool
	GetSwitchInContext context.Context
	GetSwitchInReq     *proto.GetSwitchRequest
	GetSwitchOutResp   *proto.Switch
	GetSwitchOutError  error

	GetSwitchesCalled   bool
	GetSwitchesInReq    *proto.GetSwitchesRequest
	GetSwitchesInStream proto.SwitchService_GetSwitchesServer
	GetSwitchesOutError error

	SetSwitchCalled    bool
	SetSwitchInContext context.Context
	SetSwitchInReq     *proto.SetSwitchRequest
	SetSwitchOutResp   *proto.SetSwitchResponse
	SetSwitchOutError  error
}

func (m *mockSwitchService) InstallSwitch(ctx context.Context, req *proto.InstallSwitchRequest) (*proto.Switch, error) {
	m.InstallSwitchCalled = true
	m.InstallSwitchInContext = ctx
	m.InstallSwitchInReq = req
	return m.InstallSwitchOutResp, m.InstallSwitchOutError
}

func (m *mockSwitchService) RemoveSwitch(ctx context.Context, req *proto.RemoveSwitchRequest) (*proto.RemoveSwitchResponse, error) {
	m.RemoveSwitchCalled = true
	m.RemoveSwitchInContext = ctx
	m.RemoveSwitchInReq = req
	return m.RemoveSwitchOutResp, m.RemoveSwitchOutError
}

func (m *mockSwitchService) GetSwitch(ctx context.Context, req *proto.GetSwitchRequest) (*proto.Switch, error) {
	m.GetSwitchCalled = true
	m.GetSwitchInContext = ctx
	m.GetSwitchInReq = req
	return m.GetSwitchOutResp, m.GetSwitchOutError
}

func (m *mockSwitchService) GetSwitches(req *proto.GetSwitchesRequest, stream proto.SwitchService_GetSwitchesServer) error {
	m.GetSwitchesCalled = true
	m.GetSwitchesInReq = req
	m.GetSwitchesInStream = stream
	return m.GetSwitchesOutError
}

func (m *mockSwitchService) SetSwitch(ctx context.Context, req *proto.SetSwitchRequest) (*proto.SetSwitchResponse, error) {
	m.SetSwitchCalled = true
	m.SetSwitchInContext = ctx
	m.SetSwitchInReq = req
	return m.SetSwitchOutResp, m.SetSwitchOutError
}
