package transport

import (
	"context"
	"encoding/json"

	"github.com/moorara/microservices-demo/services/asset-service/internal/queue"
	"github.com/moorara/microservices-demo/services/asset-service/internal/service"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/nats-io/go-nats"
	"github.com/opentracing/opentracing-go"

	opentracingLog "github.com/opentracing/opentracing-go/log"
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

func (t *natsTransport) extractParentSpanContext(req request) (opentracing.SpanContext, error) {
	// Get span context data
	data := make(map[string]string)
	err := json.Unmarshal([]byte(req.Span), &data)
	if err != nil {
		return nil, err
	}

	// Extract span context
	carrier := opentracing.TextMapCarrier(data)
	return t.tracer.Extract(opentracing.TextMap, carrier)
}

func (t *natsTransport) createSpan(req request) opentracing.Span {
	var span opentracing.Span
	opName := req.Kind

	parentSpanContext, _ := t.extractParentSpanContext(req)
	if parentSpanContext == nil {
		span = t.tracer.StartSpan(opName)
	} else {
		span = t.tracer.StartSpan(opName, opentracing.ChildOf(parentSpanContext))
	}

	return span
}

func (t *natsTransport) reply(subject string, response interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		t.logger.Error("message", "invalid response", "error", err)
		return
	}

	err = t.conn.Publish(subject, data)
	if err != nil {
		t.logger.Error("message", "Error publishing reply", "error", err)
		return
	}
}

func (t *natsTransport) createAlarmRequest(ctx context.Context, msg *nats.Msg) {
	var req createAlarmRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := createAlarmResponse{
		response: response{
			Kind: createAlarm,
		},
	}

	alarm, err := t.alarmService.Create(ctx, req.Input)
	res.response.Error = err
	res.Alarm = alarm

	t.reply(msg.Reply, res)
}

func (t *natsTransport) allAlarmRequest(ctx context.Context, msg *nats.Msg) {
	var req allAlarmRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := allAlarmResponse{
		response: response{
			Kind: allAlarm,
		},
	}

	alarms, err := t.alarmService.All(ctx, req.SiteID)
	res.response.Error = err
	res.Alarms = alarms

	t.reply(msg.Reply, res)
}

func (t *natsTransport) getAlarmRequest(ctx context.Context, msg *nats.Msg) {
	var req getAlarmRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := getAlarmResponse{
		response: response{
			Kind: getAlarm,
		},
	}

	alarm, err := t.alarmService.Get(ctx, req.ID)
	res.response.Error = err
	res.Alarm = alarm

	t.reply(msg.Reply, res)
}

func (t *natsTransport) updateAlarmRequest(ctx context.Context, msg *nats.Msg) {
	var req updateAlarmRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := updateAlarmResponse{
		response: response{
			Kind: updateAlarm,
		},
	}

	updated, err := t.alarmService.Update(ctx, req.ID, req.Input)
	res.response.Error = err
	res.Updated = updated

	t.reply(msg.Reply, res)
}

func (t *natsTransport) deleteAlarmRequest(ctx context.Context, msg *nats.Msg) {
	var req deleteAlarmRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := deleteAlarmResponse{
		response: response{
			Kind: deleteAlarm,
		},
	}

	deleted, err := t.alarmService.Delete(ctx, req.ID)
	res.response.Error = err
	res.Deleted = deleted

	t.reply(msg.Reply, res)
}

func (t *natsTransport) createCameraRequest(ctx context.Context, msg *nats.Msg) {
	var req createCameraRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := createCameraResponse{
		response: response{
			Kind: createCamera,
		},
	}

	camera, err := t.cameraService.Create(ctx, req.Input)
	res.response.Error = err
	res.Camera = camera

	t.reply(msg.Reply, res)
}

func (t *natsTransport) allCameraRequest(ctx context.Context, msg *nats.Msg) {
	var req allCameraRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := allCameraResponse{
		response: response{
			Kind: allCamera,
		},
	}

	cameras, err := t.cameraService.All(ctx, req.SiteID)
	res.response.Error = err
	res.Cameras = cameras

	t.reply(msg.Reply, res)
}

func (t *natsTransport) getCameraRequest(ctx context.Context, msg *nats.Msg) {
	var req getCameraRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := getCameraResponse{
		response: response{
			Kind: getCamera,
		},
	}

	camera, err := t.cameraService.Get(ctx, req.ID)
	res.response.Error = err
	res.Camera = camera

	t.reply(msg.Reply, res)
}

func (t *natsTransport) updateCameraRequest(ctx context.Context, msg *nats.Msg) {
	var req updateCameraRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := updateCameraResponse{
		response: response{
			Kind: updateCamera,
		},
	}

	updated, err := t.cameraService.Update(ctx, req.ID, req.Input)
	res.response.Error = err
	res.Updated = updated

	t.reply(msg.Reply, res)
}

func (t *natsTransport) deleteCameraRequest(ctx context.Context, msg *nats.Msg) {
	var req deleteCameraRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		t.logger.Warn("message", "invalid request", "error", err)
		return
	}

	res := deleteCameraResponse{
		response: response{
			Kind: deleteCamera,
		},
	}

	deleted, err := t.cameraService.Delete(ctx, req.ID)
	res.response.Error = err
	res.Deleted = deleted

	t.reply(msg.Reply, res)
}

func (t *natsTransport) Start() (err error) {
	t.subscription, err = t.conn.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		t.logger.Debug("message", "request received", "data", string(msg.Data))

		var req request
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			t.logger.Warn("message", "invalid request", "error", err)
			return
		}

		span := t.createSpan(req)
		span.SetTag("broker", "nats")
		span.SetTag("subject", msg.Subject)
		span.SetTag("reply", msg.Reply)
		span.LogFields(opentracingLog.String("message", string(msg.Data)))
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(context.Background(), span)

		switch req.Kind {
		case createAlarm:
			t.createAlarmRequest(ctx, msg)
		case allAlarm:
			t.allAlarmRequest(ctx, msg)
		case getAlarm:
			t.getAlarmRequest(ctx, msg)
		case updateAlarm:
			t.updateAlarmRequest(ctx, msg)
		case deleteAlarm:
			t.deleteAlarmRequest(ctx, msg)
		case createCamera:
			t.createCameraRequest(ctx, msg)
		case allCamera:
			t.allCameraRequest(ctx, msg)
		case getCamera:
			t.getCameraRequest(ctx, msg)
		case updateCamera:
			t.updateCameraRequest(ctx, msg)
		case deleteCamera:
			t.deleteCameraRequest(ctx, msg)
		default:
			t.logger.Warn("message", "unknown request", "kind", req.Kind)
		}
	})

	return err
}

func (t *natsTransport) Stop(ctx context.Context) error {
	t.subscription.Unsubscribe()
	t.conn.Close()
	return nil
}
