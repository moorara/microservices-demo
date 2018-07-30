package service

import (
	"context"

	arango "github.com/arangodb/go-driver"
	"google.golang.org/grpc/metadata"
)

// mockServerStream is a mock implementation of grpc.ServerStream
type mockServerStream struct {
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

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	m.SetHeaderCallCount++
	m.SetHeaderInMeta = md
	return m.SetHeaderOutError
}

func (m *mockServerStream) SendHeader(md metadata.MD) error {
	m.SendHeaderCallCount++
	m.SendHeaderInMeta = md
	return m.SendHeaderOutError
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
	m.SetTrailerCallCount++
	m.SetTrailerInMeta = md
}

func (m *mockServerStream) Context() context.Context {
	m.ContextCallCount++
	return m.ContextOutContext
}

func (m *mockServerStream) SendMsg(msg interface{}) error {
	m.SendMsgCallCount++
	m.SendMsgInMsg = msg
	return m.SendMsgOutError
}

func (m *mockServerStream) RecvMsg(msg interface{}) error {
	m.RecvMsgCallCount++
	m.RecvMsgInMsg = msg
	return m.RecvMsgOutError
}

// mockArangoService is a mock implementation of service.ArangoService
type mockArangoService struct {
	ConnectCallCount    int
	ConnectInContext    context.Context
	ConnectInEndpoints  []string
	ConnectInUser       string
	ConnectInPassword   string
	ConnectInDatabase   string
	ConnectInCollection string
	ConnectOutError     error

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

func (m *mockArangoService) Connect(ctx context.Context, endpoints []string, user, password, database, collection string) error {
	m.ConnectCallCount++
	m.ConnectInContext = ctx
	m.ConnectInEndpoints = endpoints
	m.ConnectInUser = user
	m.ConnectInPassword = password
	m.ConnectInDatabase = database
	m.ConnectInCollection = collection
	return m.ConnectOutError
}

func (m *mockArangoService) Query(ctx context.Context, query string, vars map[string]interface{}) (arango.Cursor, error) {
	m.QueryCallCount++
	m.QueryInContext = ctx
	m.QueryInQuery = query
	m.QueryInVars = vars
	return m.QueryOutCursor, m.QueryOutError
}

func (m *mockArangoService) CreateDocument(ctx context.Context, doc interface{}) (arango.DocumentMeta, error) {
	m.CreateDocumentCallCount++
	m.CreateDocumentInContext = ctx
	m.CreateDocumentInDoc = doc
	return m.CreateDocumentOutMeta, m.CreateDocumentOutError
}

func (m *mockArangoService) ReadDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.ReadDocumentCallCount++
	m.ReadDocumentInContext = ctx
	m.ReadDocumentInKey = key
	m.ReadDocumentInDoc = doc
	return m.ReadDocumentOutMeta, m.ReadDocumentOutError
}

func (m *mockArangoService) UpdateDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.UpdateDocumentCallCount++
	m.UpdateDocumentInContext = ctx
	m.UpdateDocumentInKey = key
	m.UpdateDocumentInDoc = doc
	return m.UpdateDocumentOutMeta, m.UpdateDocumentOutError
}

func (m *mockArangoService) RemoveDocument(ctx context.Context, key string) (arango.DocumentMeta, error) {
	m.RemoveDocumentCallCount++
	m.RemoveDocumentInContext = ctx
	m.RemoveDocumentInKey = key
	return m.RemoveDocumentOutMeta, m.RemoveDocumentOutError
}
