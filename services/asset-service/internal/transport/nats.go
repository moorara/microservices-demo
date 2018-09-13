package transport

import (
	"context"
	"encoding/json"

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
		Subscribe() error
	}

	natsTransport struct {
		logger        *log.Logger
		metrics       *metrics.Metrics
		tracer        opentracing.Tracer
		conn          queue.NATSConnection
		alarmService  service.AlarmService
		cameraService service.CameraService
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

func (t *natsTransport) Subscribe() error {
	_, err := t.conn.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		var req request
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			t.handlerError(err)
			return
		}

		switch req.Kind {
		case createAlarm:
			t.createAlarm(msg.Data, msg.Reply)
		case allAlarm:
			t.allAlarm(msg.Data, msg.Reply)
		case getAlarm:
			t.getAlarm(msg.Data, msg.Reply)
		case updateAlarm:
			t.updateAlarm(msg.Data, msg.Reply)
		case deleteAlarm:
			t.deleteAlarm(msg.Data, msg.Reply)
		case createCamera:
			t.createCamera(msg.Data, msg.Reply)
		case allCamera:
			t.allCamera(msg.Data, msg.Reply)
		case getCamera:
			t.getCamera(msg.Data, msg.Reply)
		case updateCamera:
			t.updateCamera(msg.Data, msg.Reply)
		case deleteCamera:
			t.deleteCamera(msg.Data, msg.Reply)
		}
	})

	return err
}

func (t *natsTransport) handlerError(err error) {

}

func (t *natsTransport) createContext(req request) context.Context {
	return context.Background()
}

func (t *natsTransport) replyResponse(reply string, data []byte) {

}

func (t *natsTransport) createAlarm(data []byte, reply string) {
	var req createAlarmRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}

	ctx := t.createContext(req.request)
	alarm, err := t.alarmService.Create(ctx, req.input)
	if err != nil {
		t.handlerError(err)
	}

	res := createAlarmResponse{
		response: response{
			Kind: req.request.Kind,
			Span: req.request.Span,
		},
		Alarm: *alarm,
		Error: err,
	}

	out, err := json.Marshal(&res)
	if err != nil {
		t.handlerError(err)
	}

	t.replyResponse(reply, out)
}

func (t *natsTransport) allAlarm(data []byte, reply string) {
	var req allAlarmRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) getAlarm(data []byte, reply string) {
	var req getAlarmRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) updateAlarm(data []byte, reply string) {
	var req updateAlarmRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) deleteAlarm(data []byte, reply string) {
	var req deleteAlarmRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) createCamera(data []byte, reply string) {
	var req createCameraRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) allCamera(data []byte, reply string) {
	var req allCameraRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) getCamera(data []byte, reply string) {
	var req getCameraRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) updateCamera(data []byte, reply string) {
	var req updateCameraRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}

func (t *natsTransport) deleteCamera(data []byte, reply string) {
	var req deleteCameraRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		t.handlerError(err)
	}
}
