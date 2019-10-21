package transport

import (
	"context"
	"time"

	"github.com/moorara/microservices-demo/services/asset/internal/model"
	nats "github.com/nats-io/go-nats"
)

type mockNATSConnection struct {
	CloseCalled bool

	FlushCalled   bool
	FlushOutError error

	LastErrorCalled   bool
	LastErrorOutError error

	PublishCalled    bool
	PublishInSubject string
	PublishInData    []byte
	PublishOutError  error

	PublishMsgCalled   bool
	PublishMsgInMsg    *nats.Msg
	PublishMsgOutError error

	PublishRequestCalled    bool
	PublishRequestInSubject string
	PublishRequestInReply   string
	PublishRequestInData    []byte
	PublishRequestOutError  error

	QueueSubscribeCalled          bool
	QueueSubscribeInSubject       string
	QueueSubscribeInQueue         string
	QueueSubscribeInCallback      nats.MsgHandler
	QueueSubscribeOutSubscription *nats.Subscription
	QueueSubscribeOutError        error

	QueueSubscribeSyncCalled          bool
	QueueSubscribeSyncInSubject       string
	QueueSubscribeSyncInQueue         string
	QueueSubscribeSyncOutSubscription *nats.Subscription
	QueueSubscribeSyncOutError        error

	QueueSubscribeSyncWithChanCalled          bool
	QueueSubscribeSyncWithChanInSubject       string
	QueueSubscribeSyncWithChanInQueue         string
	QueueSubscribeSyncWithChanInChannel       chan *nats.Msg
	QueueSubscribeSyncWithChanOutSubscription *nats.Subscription
	QueueSubscribeSyncWithChanOutError        error

	RequestCalled    bool
	RequestInSubject string
	RequestInData    []byte
	RequestInTimeout time.Duration
	RequestOutMsg    *nats.Msg
	RequestOutError  error

	RequestWithContextCalled    bool
	RequestWithContextInContext context.Context
	RequestWithContextInSubject string
	RequestWithContextInData    []byte
	RequestWithContextOutMsg    *nats.Msg
	RequestWithContextOutError  error

	SubscribeCalled          bool
	SubscribeInSubject       string
	SubscribeInCallback      nats.MsgHandler
	SubscribeOutSubscription *nats.Subscription
	SubscribeOutError        error

	SubscribeSyncCalled          bool
	SubscribeSyncInSubject       string
	SubscribeSyncOutSubscription *nats.Subscription
	SubscribeSyncOutError        error
}

func (m *mockNATSConnection) Close() {
	m.CloseCalled = true
}

func (m *mockNATSConnection) Flush() error {
	m.FlushCalled = true
	return m.FlushOutError
}

func (m *mockNATSConnection) LastError() error {
	m.LastErrorCalled = true
	return m.LastErrorOutError
}

func (m *mockNATSConnection) Publish(subject string, data []byte) error {
	m.PublishCalled = true
	m.PublishInSubject = subject
	m.PublishInData = data
	return m.PublishOutError
}

func (m *mockNATSConnection) PublishMsg(msg *nats.Msg) error {
	m.PublishMsgCalled = true
	m.PublishMsgInMsg = msg
	return m.PublishMsgOutError
}

func (m *mockNATSConnection) PublishRequest(subject, reply string, data []byte) error {
	m.PublishRequestCalled = true
	m.PublishRequestInSubject = subject
	m.PublishRequestInReply = reply
	m.PublishRequestInData = data
	return m.PublishRequestOutError
}

func (m *mockNATSConnection) QueueSubscribe(subject, queue string, callback nats.MsgHandler) (*nats.Subscription, error) {
	m.QueueSubscribeCalled = true
	m.QueueSubscribeInSubject = subject
	m.QueueSubscribeInQueue = queue
	m.QueueSubscribeInCallback = callback
	return m.QueueSubscribeOutSubscription, m.QueueSubscribeOutError
}

func (m *mockNATSConnection) QueueSubscribeSync(subject, queue string) (*nats.Subscription, error) {
	m.QueueSubscribeSyncCalled = true
	m.QueueSubscribeSyncInSubject = subject
	m.QueueSubscribeSyncInQueue = queue
	return m.QueueSubscribeSyncOutSubscription, m.QueueSubscribeSyncOutError
}

func (m *mockNATSConnection) QueueSubscribeSyncWithChan(subject, queue string, channel chan *nats.Msg) (*nats.Subscription, error) {
	m.QueueSubscribeSyncWithChanCalled = true
	m.QueueSubscribeSyncWithChanInSubject = subject
	m.QueueSubscribeSyncWithChanInQueue = queue
	m.QueueSubscribeSyncWithChanInChannel = channel
	return m.QueueSubscribeSyncWithChanOutSubscription, m.QueueSubscribeSyncWithChanOutError
}

func (m *mockNATSConnection) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	m.RequestCalled = true
	m.RequestInSubject = subject
	m.RequestInData = data
	m.RequestInTimeout = timeout
	return m.RequestOutMsg, m.RequestOutError
}

func (m *mockNATSConnection) RequestWithContext(ctx context.Context, subject string, data []byte) (*nats.Msg, error) {
	m.RequestWithContextCalled = true
	m.RequestWithContextInContext = ctx
	m.RequestWithContextInSubject = subject
	m.RequestWithContextInData = data
	return m.RequestWithContextOutMsg, m.RequestWithContextOutError
}

func (m *mockNATSConnection) Subscribe(subject string, callback nats.MsgHandler) (*nats.Subscription, error) {
	m.SubscribeCalled = true
	m.SubscribeInSubject = subject
	m.SubscribeInCallback = callback
	return m.SubscribeOutSubscription, m.SubscribeOutError
}

func (m *mockNATSConnection) SubscribeSync(subject string) (*nats.Subscription, error) {
	m.SubscribeSyncCalled = true
	m.SubscribeSyncInSubject = subject
	return m.SubscribeSyncOutSubscription, m.SubscribeSyncOutError
}

type mockAlarmService struct {
	CreateCalled    bool
	CreateInContext context.Context
	CreateInInput   model.AlarmInput
	CreateOutAlarm  *model.Alarm
	CreateOutError  error

	AllCalled    bool
	AllInContext context.Context
	AllInSiteID  string
	AllOutAlarms []model.Alarm
	AllOutError  error

	GetCalled    bool
	GetInContext context.Context
	GetInID      string
	GetOutAlarm  *model.Alarm
	GetOutError  error

	UpdateCalled     bool
	UpdateInContext  context.Context
	UpdateInID       string
	UpdateInInput    model.AlarmInput
	UpdateOutUpdated bool
	UpdateOutError   error

	DeleteCalled     bool
	DeleteInContext  context.Context
	DeleteInID       string
	DeleteOutDeleted bool
	DeleteOutError   error
}

func (m *mockAlarmService) Create(ctx context.Context, input model.AlarmInput) (*model.Alarm, error) {
	m.CreateCalled = true
	m.CreateInContext = ctx
	m.CreateInInput = input
	return m.CreateOutAlarm, m.CreateOutError
}

func (m *mockAlarmService) All(ctx context.Context, siteID string) ([]model.Alarm, error) {
	m.AllCalled = true
	m.AllInContext = ctx
	m.AllInSiteID = siteID
	return m.AllOutAlarms, m.AllOutError
}

func (m *mockAlarmService) Get(ctx context.Context, id string) (*model.Alarm, error) {
	m.GetCalled = true
	m.GetInContext = ctx
	m.GetInID = id
	return m.GetOutAlarm, m.GetOutError
}

func (m *mockAlarmService) Update(ctx context.Context, id string, input model.AlarmInput) (bool, error) {
	m.UpdateCalled = true
	m.UpdateInContext = ctx
	m.UpdateInID = id
	m.UpdateInInput = input
	return m.UpdateOutUpdated, m.UpdateOutError
}

func (m *mockAlarmService) Delete(ctx context.Context, id string) (bool, error) {
	m.DeleteCalled = true
	m.DeleteInContext = ctx
	m.DeleteInID = id
	return m.DeleteOutDeleted, m.DeleteOutError
}

type mockCameraService struct {
	CreateCalled    bool
	CreateInContext context.Context
	CreateInInput   model.CameraInput
	CreateOutCamera *model.Camera
	CreateOutError  error

	AllCalled     bool
	AllInContext  context.Context
	AllInSiteID   string
	AllOutCameras []model.Camera
	AllOutError   error

	GetCalled    bool
	GetInContext context.Context
	GetInID      string
	GetOutCamera *model.Camera
	GetOutError  error

	UpdateCalled     bool
	UpdateInContext  context.Context
	UpdateInID       string
	UpdateInInput    model.CameraInput
	UpdateOutUpdated bool
	UpdateOutError   error

	DeleteCalled     bool
	DeleteInContext  context.Context
	DeleteInID       string
	DeleteOutDeleted bool
	DeleteOutError   error
}

func (m *mockCameraService) Create(ctx context.Context, input model.CameraInput) (*model.Camera, error) {
	m.CreateCalled = true
	m.CreateInContext = ctx
	m.CreateInInput = input
	return m.CreateOutCamera, m.CreateOutError
}

func (m *mockCameraService) All(ctx context.Context, siteID string) ([]model.Camera, error) {
	m.AllCalled = true
	m.AllInContext = ctx
	m.AllInSiteID = siteID
	return m.AllOutCameras, m.AllOutError
}

func (m *mockCameraService) Get(ctx context.Context, id string) (*model.Camera, error) {
	m.GetCalled = true
	m.GetInContext = ctx
	m.GetInID = id
	return m.GetOutCamera, m.GetOutError
}

func (m *mockCameraService) Update(ctx context.Context, id string, input model.CameraInput) (bool, error) {
	m.UpdateCalled = true
	m.UpdateInContext = ctx
	m.UpdateInID = id
	m.UpdateInInput = input
	return m.UpdateOutUpdated, m.UpdateOutError
}

func (m *mockCameraService) Delete(ctx context.Context, id string) (bool, error) {
	m.DeleteCalled = true
	m.DeleteInContext = ctx
	m.DeleteInID = id
	return m.DeleteOutDeleted, m.DeleteOutError
}
