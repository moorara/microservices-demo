package service

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type mockServerStream struct {
	SetHeaderCallCount    int
	SetHeaderInMetatadata metadata.MD
	SetHeaderOutError     error

	SendHeaderCallCount    int
	SendHeaderInMetatadata metadata.MD
	SendHeaderOutError     error

	SetTrailerCallCount    int
	SetTrailerInMetatadata metadata.MD

	ContextCallCount  int
	ContextOutContext context.Context

	SendMsgCallCount int
	SendMsgInMessage interface{}
	SendMsgOutError  error

	RecvMsgCallCount int
	RecvMsgInMessage interface{}
	RecvMsgOutError  error
}

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	m.SetHeaderCallCount++
	m.SetHeaderInMetatadata = md
	return m.SetHeaderOutError
}

func (m *mockServerStream) SendHeader(md metadata.MD) error {
	m.SendHeaderCallCount++
	m.SendHeaderInMetatadata = md
	return m.SendHeaderOutError
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
	m.SetTrailerCallCount++
	m.SetTrailerInMetatadata = md
}

func (m *mockServerStream) Context() context.Context {
	m.ContextCallCount++
	return m.ContextOutContext
}

func (m *mockServerStream) SendMsg(msg interface{}) error {
	m.SendMsgCallCount++
	m.SendMsgInMessage = msg
	return m.SendMsgOutError
}

func (m *mockServerStream) RecvMsg(msg interface{}) error {
	m.RecvMsgCallCount++
	m.RecvMsgInMessage = msg
	return m.RecvMsgOutError
}
