package event

import (
	"context"
	"time"

	nats "github.com/nats-io/go-nats"
)

const (
	natsMaxReconnect  = 10
	natsReconnectWait = 2 * time.Second
)

type (
	// NATSConnection is the connection to NATS cluster
	NATSConnection interface {
		Close()
		Flush() error
		LastError() error
		Publish(subject string, data []byte) error
		PublishMsg(msg *nats.Msg) error
		PublishRequest(subject, reply string, data []byte) error
		QueueSubscribe(subject, queue string, callback nats.MsgHandler) (*nats.Subscription, error)
		QueueSubscribeSync(subject, queue string) (*nats.Subscription, error)
		QueueSubscribeSyncWithChan(subject, queue string, channel chan *nats.Msg) (*nats.Subscription, error)
		Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error)
		RequestWithContext(ctx context.Context, subject string, data []byte) (*nats.Msg, error)
		Subscribe(subject string, callback nats.MsgHandler) (*nats.Subscription, error)
		SubscribeSync(subject string) (*nats.Subscription, error)
	}
)

// NewNATSConnection creates a new connection to a NATS cluster
func NewNATSConnection(servers []string, name, user, password string) (NATSConnection, error) {
	opts := nats.Options{
		Servers:        servers,
		Name:           name,
		User:           user,
		Password:       password,
		AllowReconnect: true,
		MaxReconnect:   natsMaxReconnect,
		ReconnectWait:  natsReconnectWait,
	}

	return opts.Connect()
}
