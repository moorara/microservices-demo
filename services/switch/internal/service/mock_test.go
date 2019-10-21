package service

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	arango "github.com/arangodb/go-driver"
	"github.com/moorara/microservices-demo/services/switch/internal/proto"
)

// mockCloser is a mock implementation of io.Closer
type mockCloser struct {
	CloseCalled   bool
	CloseOutError error
}

func (m *mockCloser) Close() error {
	m.CloseCalled = true
	return m.CloseOutError
}

// mockServerStream is a mock implementation of grpc.ServerStream
type mockServerStream struct {
	SetHeaderCalled   bool
	SetHeaderInMeta   metadata.MD
	SetHeaderOutError error

	SendHeaderCalled   bool
	SendHeaderInMeta   metadata.MD
	SendHeaderOutError error

	SetTrailerCalled bool
	SetTrailerInMeta metadata.MD

	ContextCalled     bool
	ContextOutContext context.Context

	SendMsgCalled   bool
	SendMsgInMsg    interface{}
	SendMsgOutError error

	RecvMsgCalled   bool
	RecvMsgInMsg    interface{}
	RecvMsgOutError error
}

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	m.SetHeaderCalled = true
	m.SetHeaderInMeta = md
	return m.SetHeaderOutError
}

func (m *mockServerStream) SendHeader(md metadata.MD) error {
	m.SendHeaderCalled = true
	m.SendHeaderInMeta = md
	return m.SendHeaderOutError
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
	m.SetTrailerCalled = true
	m.SetTrailerInMeta = md
}

func (m *mockServerStream) Context() context.Context {
	m.ContextCalled = true
	return m.ContextOutContext
}

func (m *mockServerStream) SendMsg(msg interface{}) error {
	m.SendMsgCalled = true
	m.SendMsgInMsg = msg
	return m.SendMsgOutError
}

func (m *mockServerStream) RecvMsg(msg interface{}) error {
	m.RecvMsgCalled = true
	m.RecvMsgInMsg = msg
	return m.RecvMsgOutError
}

// mockGetSwitchesServer mocks proto.SwitchService_GetSwitchesServer
type mockGetSwitchesServer struct {
	grpc.ServerStream

	SendCalled   bool
	SendInSwitch *proto.Switch
	SendOutError error
}

func (m *mockGetSwitchesServer) Send(sw *proto.Switch) error {
	m.SendCalled = true
	m.SendInSwitch = sw
	return m.SendOutError
}

// mockArangoCursor is a mock implementation of arango.Cursor
type mockArangoCursor struct {
	io.Closer

	HasMoreCallCount  int
	HasMoreOutResults []bool

	ReadDocumentCalled    bool
	ReadDocumentInContext context.Context
	ReadDocumentInDoc     interface{}
	ReadDocumentOutMeta   arango.DocumentMeta
	ReadDocumentOutError  error

	CountCalled    bool
	CountOutResult int64

	StatisticsCalled    bool
	StatisticsOutResult arango.QueryStatistics
}

func (m *mockArangoCursor) HasMore() bool {
	i := m.HasMoreCallCount % len(m.HasMoreOutResults)
	m.HasMoreCallCount++
	return m.HasMoreOutResults[i]
}

func (m *mockArangoCursor) ReadDocument(ctx context.Context, doc interface{}) (arango.DocumentMeta, error) {
	m.ReadDocumentCalled = true
	m.ReadDocumentInContext = ctx
	m.ReadDocumentInDoc = doc
	return m.ReadDocumentOutMeta, m.ReadDocumentOutError
}

func (m *mockArangoCursor) Count() int64 {
	m.CountCalled = true
	return m.CountOutResult
}

func (m *mockArangoCursor) Statistics() arango.QueryStatistics {
	m.StatisticsCalled = true
	return m.StatisticsOutResult
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

	CreateDocumentCalled    bool
	CreateDocumentInContext context.Context
	CreateDocumentInDoc     interface{}
	CreateDocumentOutMeta   arango.DocumentMeta
	CreateDocumentOutError  error

	ReadDocumentCalled    bool
	ReadDocumentInContext context.Context
	ReadDocumentInKey     string
	ReadDocumentInDoc     interface{}
	ReadDocumentOutMeta   arango.DocumentMeta
	ReadDocumentOutError  error

	UpdateDocumentCalled    bool
	UpdateDocumentInContext context.Context
	UpdateDocumentInKey     string
	UpdateDocumentInDoc     interface{}
	UpdateDocumentOutMeta   arango.DocumentMeta
	UpdateDocumentOutError  error

	RemoveDocumentCalled    bool
	RemoveDocumentInContext context.Context
	RemoveDocumentInKey     string
	RemoveDocumentOutMeta   arango.DocumentMeta
	RemoveDocumentOutError  error
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
	return m.CreateDocumentOutMeta, m.CreateDocumentOutError
}

func (m *mockArangoService) ReadDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.ReadDocumentCalled = true
	m.ReadDocumentInContext = ctx
	m.ReadDocumentInKey = key
	m.ReadDocumentInDoc = doc
	return m.ReadDocumentOutMeta, m.ReadDocumentOutError
}

func (m *mockArangoService) UpdateDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.UpdateDocumentCalled = true
	m.UpdateDocumentInContext = ctx
	m.UpdateDocumentInKey = key
	m.UpdateDocumentInDoc = doc
	return m.UpdateDocumentOutMeta, m.UpdateDocumentOutError
}

func (m *mockArangoService) RemoveDocument(ctx context.Context, key string) (arango.DocumentMeta, error) {
	m.RemoveDocumentCalled = true
	m.RemoveDocumentInContext = ctx
	m.RemoveDocumentInKey = key
	return m.RemoveDocumentOutMeta, m.RemoveDocumentOutError
}
