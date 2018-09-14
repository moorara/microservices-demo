package transport

import (
	"context"

	nats "github.com/nats-io/go-nats"

	"github.com/moorara/microservices-demo/services/asset-service/internal/queue"
	"github.com/moorara/microservices-demo/services/asset-service/internal/service"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/opentracing/opentracing-go"
)

const (
	subject    = "asset_service"
	queueGroup = "workers"
)

type (
	// NATSTransport is the transport for NATS
	NATSTransport interface {
		Start() error
		Stop(context.Context) error
	}

	natsTransport struct {
		logger        *log.Logger
		metrics       *metrics.Metrics
		tracer        opentracing.Tracer
		conn          queue.NATSConnection
		alarmService  service.AlarmService
		cameraService service.CameraService
		subscription  *nats.Subscription
	}
)

// NewNATSTransport creates a new NATS transport instance
func NewNATSTransport(logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer,
	conn queue.NATSConnection, alarmService service.AlarmService, cameraService service.CameraService) NATSTransport {
	return &natsTransport{
		logger:        logger,
		metrics:       metrics,
		tracer:        tracer,
		conn:          conn,
		alarmService:  alarmService,
		cameraService: cameraService,
	}
}

func (t *natsTransport) Start() (err error) {
	t.subscription, err = t.conn.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		t.logger.Warn("message", string(msg.Data))
		t.conn.Publish(msg.Reply, []byte("Your message is received!"))
	})

	return err
}

func (t *natsTransport) Stop(ctx context.Context) error {
	t.subscription.Unsubscribe()
	t.conn.Close()
	return nil
}
