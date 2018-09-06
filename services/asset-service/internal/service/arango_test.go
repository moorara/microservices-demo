package service

import (
	"context"
	"io"

	arango "github.com/arangodb/go-driver"
)

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

// mockArango is a mock implementation of db.Arango
type mockArango struct {
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

func (m *mockArango) Connect(ctx context.Context, database, collection string) error {
	m.ConnectCalled = true
	m.ConnectInContext = ctx
	m.ConnectInDatabase = database
	m.ConnectInCollection = collection
	return m.ConnectOutError
}

func (m *mockArango) Query(ctx context.Context, query string, vars map[string]interface{}) (arango.Cursor, error) {
	m.QueryCalled = true
	m.QueryInContext = ctx
	m.QueryInQuery = query
	m.QueryInVars = vars
	return m.QueryOutCursor, m.QueryOutError
}

func (m *mockArango) CreateDocument(ctx context.Context, doc interface{}) (arango.DocumentMeta, error) {
	m.CreateDocumentCalled = true
	m.CreateDocumentInContext = ctx
	m.CreateDocumentInDoc = doc
	return m.CreateDocumentOutMeta, m.CreateDocumentOutError
}

func (m *mockArango) ReadDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.ReadDocumentCalled = true
	m.ReadDocumentInContext = ctx
	m.ReadDocumentInKey = key
	m.ReadDocumentInDoc = doc
	return m.ReadDocumentOutMeta, m.ReadDocumentOutError
}

func (m *mockArango) UpdateDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error) {
	m.UpdateDocumentCalled = true
	m.UpdateDocumentInContext = ctx
	m.UpdateDocumentInKey = key
	m.UpdateDocumentInDoc = doc
	return m.UpdateDocumentOutMeta, m.UpdateDocumentOutError
}

func (m *mockArango) RemoveDocument(ctx context.Context, key string) (arango.DocumentMeta, error) {
	m.RemoveDocumentCalled = true
	m.RemoveDocumentInContext = ctx
	m.RemoveDocumentInKey = key
	return m.RemoveDocumentOutMeta, m.RemoveDocumentOutError
}
