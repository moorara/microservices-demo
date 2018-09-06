package service

import (
	"context"
	"time"

	nats "github.com/nats-io/go-nats"
)

// mockArango is a mock implementation of event.NATSConnection
type mockNATSConnection struct {
	CloseCallCount int

	FlushCallCount int
	FlushOutError  error

	LastErrorCallCount int
	LastErrorOutError  error

	PublishCallCount int
	PublishInSubject string
	PublishInData    []byte
	PublishOutError  error

	PublishMsgCallCount int
	PublishMsgInMsg     *nats.Msg
	PublishMsgOutError  error

	PublishRequestCallCount int
	PublishRequestInSubject string
	PublishRequestInReply   string
	PublishRequestInData    []byte
	PublishRequestOutError  error

	QueueSubscribeCallCount       int
	QueueSubscribeInSubject       string
	QueueSubscribeInQueue         string
	QueueSubscribeInCallback      nats.MsgHandler
	QueueSubscribeOutSubscription *nats.Subscription
	QueueSubscribeOutError        error

	QueueSubscribeSyncCallCount       int
	QueueSubscribeSyncInSubject       string
	QueueSubscribeSyncInQueue         string
	QueueSubscribeSyncOutSubscription *nats.Subscription
	QueueSubscribeSyncOutError        error

	QueueSubscribeSyncWithChanCallCount       int
	QueueSubscribeSyncWithChanInSubject       string
	QueueSubscribeSyncWithChanInQueue         string
	QueueSubscribeSyncWithChanInChannel       chan *nats.Msg
	QueueSubscribeSyncWithChanOutSubscription *nats.Subscription
	QueueSubscribeSyncWithChanOutError        error

	RequestCallCount int
	RequestInSubject string
	RequestInData    []byte
	RequestInTimeout time.Duration
	RequestOutMsg    *nats.Msg
	RequestOutError  error

	RequestWithContextCallCount int
	RequestWithContextInContext context.Context
	RequestWithContextInSubject string
	RequestWithContextInData    []byte
	RequestWithContextOutMsg    *nats.Msg
	RequestWithContextOutError  error

	SubscribeCallCount       int
	SubscribeInSubject       string
	SubscribeInCallback      nats.MsgHandler
	SubscribeOutSubscription *nats.Subscription
	SubscribeOutError        error

	SubscribeSyncCallCount       int
	SubscribeSyncInSubject       string
	SubscribeSyncOutSubscription *nats.Subscription
	SubscribeSyncOutError        error
}

func (m *mockNATSConnection) Close() {
	m.CloseCallCount++
}

func (m *mockNATSConnection) Flush() error {
	m.FlushCallCount++
	return m.FlushOutError
}

func (m *mockNATSConnection) LastError() error {
	m.LastErrorCallCount++
	return m.LastErrorOutError
}

func (m *mockNATSConnection) Publish(subject string, data []byte) error {
	m.PublishCallCount++
	m.PublishInSubject = subject
	m.PublishInData = data
	return m.PublishOutError
}

func (m *mockNATSConnection) PublishMsg(msg *nats.Msg) error {
	m.PublishMsgCallCount++
	m.PublishMsgInMsg = msg
	return m.PublishMsgOutError
}

func (m *mockNATSConnection) PublishRequest(subject, reply string, data []byte) error {
	m.PublishRequestCallCount++
	m.PublishRequestInSubject = subject
	m.PublishRequestInReply = reply
	m.PublishRequestInData = data
	return m.PublishRequestOutError
}

func (m *mockNATSConnection) QueueSubscribe(subject, queue string, callback nats.MsgHandler) (*nats.Subscription, error) {
	m.QueueSubscribeCallCount++
	m.QueueSubscribeInSubject = subject
	m.QueueSubscribeInQueue = queue
	m.QueueSubscribeInCallback = callback
	return m.QueueSubscribeOutSubscription, m.QueueSubscribeOutError
}

func (m *mockNATSConnection) QueueSubscribeSync(subject, queue string) (*nats.Subscription, error) {
	m.QueueSubscribeSyncCallCount++
	m.QueueSubscribeSyncInSubject = subject
	m.QueueSubscribeSyncInQueue = queue
	return m.QueueSubscribeSyncOutSubscription, m.QueueSubscribeSyncOutError
}

func (m *mockNATSConnection) QueueSubscribeSyncWithChan(subject, queue string, channel chan *nats.Msg) (*nats.Subscription, error) {
	m.QueueSubscribeSyncWithChanCallCount++
	m.QueueSubscribeSyncWithChanInSubject = subject
	m.QueueSubscribeSyncWithChanInQueue = queue
	m.QueueSubscribeSyncWithChanInChannel = channel
	return m.QueueSubscribeSyncWithChanOutSubscription, m.QueueSubscribeSyncWithChanOutError
}

func (m *mockNATSConnection) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	m.RequestCallCount++
	m.RequestInSubject = subject
	m.RequestInData = data
	m.RequestInTimeout = timeout
	return m.RequestOutMsg, m.RequestOutError
}

func (m *mockNATSConnection) RequestWithContext(ctx context.Context, subject string, data []byte) (*nats.Msg, error) {
	m.RequestWithContextCallCount++
	m.RequestWithContextInContext = ctx
	m.RequestWithContextInSubject = subject
	m.RequestWithContextInData = data
	return m.RequestWithContextOutMsg, m.RequestWithContextOutError
}

func (m *mockNATSConnection) Subscribe(subject string, callback nats.MsgHandler) (*nats.Subscription, error) {
	m.SubscribeCallCount++
	m.SubscribeInSubject = subject
	m.SubscribeInCallback = callback
	return m.SubscribeOutSubscription, m.SubscribeOutError
}

func (m *mockNATSConnection) SubscribeSync(subject string) (*nats.Subscription, error) {
	m.SubscribeSyncCallCount++
	m.SubscribeSyncInSubject = subject
	return m.SubscribeSyncOutSubscription, m.SubscribeSyncOutError
}
