package transport

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/moorara/microservices-demo/services/asset/internal/model"
	"github.com/moorara/microservices-demo/services/asset/pkg/log"
	"github.com/moorara/microservices-demo/services/asset/pkg/metrics"
	"github.com/nats-io/go-nats"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestNewNATSTransport(t *testing.T) {
	logger := log.NewNopLogger()
	metrics := metrics.New("unit-test")
	tracer := mocktracer.New()

	conn := &mockNATSConnection{}
	alarmService := &mockAlarmService{}
	cameraService := &mockCameraService{}

	natsTransport := NewNATSTransport(logger, metrics, tracer, conn, alarmService, cameraService)
	assert.NotNil(t, natsTransport)
}

func TestExtractParentSpanContext(t *testing.T) {
	logger := log.NewNopLogger()
	metrics := metrics.New("unit-test")
	tracer := mocktracer.New()

	tests := []struct {
		name    string
		reqKind string
		reqSpan opentracing.Span
	}{
		{
			"WithoutSpan",
			createAlarm,
			nil,
		},
		{
			"WithSpan",
			createCamera,
			tracer.StartSpan("createCamera"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := request{
				Kind: tc.reqKind,
			}

			// Inject the span if any
			if tc.reqSpan != nil {
				carrier := opentracing.TextMapCarrier{}
				err := tracer.Inject(tc.reqSpan.Context(), opentracing.TextMap, carrier)
				assert.NoError(t, err)

				data, err := json.Marshal(carrier)
				assert.NoError(t, err)

				req.Span = string(data)
			}

			nt := &natsTransport{
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			parentSpanContext, err := nt.extractParentSpanContext(req)

			if tc.reqSpan == nil {
				assert.Error(t, err)
				assert.Nil(t, parentSpanContext)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parentSpanContext)
			}
		})
	}
}

func TestCreateSpan(t *testing.T) {
	logger := log.NewNopLogger()
	metrics := metrics.New("unit-test")
	tracer := mocktracer.New()

	tests := []struct {
		name    string
		reqKind string
		reqSpan opentracing.Span
	}{
		{
			"WithoutSpan",
			createAlarm,
			nil,
		},
		{
			"WithSpan",
			createCamera,
			tracer.StartSpan("createCamera"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := request{
				Kind: tc.reqKind,
			}

			// Inject the span if any
			if tc.reqSpan != nil {
				carrier := opentracing.TextMapCarrier{}
				err := tracer.Inject(tc.reqSpan.Context(), opentracing.TextMap, carrier)
				assert.NoError(t, err)

				data, err := json.Marshal(carrier)
				assert.NoError(t, err)

				req.Span = string(data)
			}

			nt := &natsTransport{
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
			}

			span := nt.createSpan(req)

			assert.NotNil(t, span)
		})
	}
}

func TestReply(t *testing.T) {
	tests := []struct {
		name         string
		subject      string
		response     interface{}
		publishError error
	}{
		{
			"NoResponse",
			"reply_subject",
			nil,
			nil,
		},
		{
			"PublishError",
			"reply_subject",
			map[string]interface{}{
				"kind": "createAlarm",
			},
			errors.New("nats publish error"),
		},
		{
			"Successful",
			"reply_subject",
			map[string]interface{}{
				"kind": "createCamera",
			},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()

			conn := &mockNATSConnection{
				PublishOutError: tc.publishError,
			}

			nt := &natsTransport{
				logger:  logger,
				metrics: metrics,
				tracer:  tracer,
				conn:    conn,
			}

			nt.reply(tc.subject, tc.response)

			assert.True(t, conn.PublishCalled)
			assert.Equal(t, tc.subject, conn.PublishInSubject)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name             string
		conn             *mockNATSConnection
		alarmService     *mockAlarmService
		cameraService    *mockCameraService
		request          map[string]interface{}
		expectedResponse map[string]interface{}
	}{
		{
			"Default",
			&mockNATSConnection{},
			&mockAlarmService{},
			&mockCameraService{},
			map[string]interface{}{},
			nil,
		},
		{
			"CreateAlarm",
			&mockNATSConnection{},
			&mockAlarmService{
				CreateOutAlarm: &model.Alarm{
					Asset: model.Asset{
						ID:       "aaaa-aaaa",
						SiteID:   "1111-1111",
						SerialNo: "1001",
					},
					Material: "co",
				},
			},
			&mockCameraService{},
			map[string]interface{}{
				"kind": createAlarm,
				"input": map[string]interface{}{
					"siteId":   "1111-1111",
					"serialNo": "1001",
					"material": "co",
				},
			},
			map[string]interface{}{
				"kind": createAlarm,
				"alarm": map[string]interface{}{
					"id":       "aaaa-aaaa",
					"siteId":   "1111-1111",
					"serialNo": "1001",
					"material": "co",
				},
			},
		},
		{
			"AllAlarm",
			&mockNATSConnection{},
			&mockAlarmService{
				AllOutAlarms: []model.Alarm{
					model.Alarm{
						Asset: model.Asset{
							ID:       "aaaa-aaaa",
							SiteID:   "1111-1111",
							SerialNo: "1001",
						},
						Material: "co",
					},
				},
			},
			&mockCameraService{},
			map[string]interface{}{
				"kind":   allAlarm,
				"siteId": "1111-1111",
			},
			map[string]interface{}{
				"kind": allAlarm,
				"alarms": []interface{}{
					map[string]interface{}{
						"id":       "aaaa-aaaa",
						"siteId":   "1111-1111",
						"serialNo": "1001",
						"material": "co",
					},
				},
			},
		},
		{
			"GetAlarm",
			&mockNATSConnection{},
			&mockAlarmService{
				GetOutAlarm: &model.Alarm{
					Asset: model.Asset{
						ID:       "aaaa-aaaa",
						SiteID:   "1111-1111",
						SerialNo: "1001",
					},
					Material: "co",
				},
			},
			&mockCameraService{},
			map[string]interface{}{
				"kind": getAlarm,
				"id":   "aaaa-aaaa",
			},
			map[string]interface{}{
				"kind": getAlarm,
				"alarm": map[string]interface{}{
					"id":       "aaaa-aaaa",
					"siteId":   "1111-1111",
					"serialNo": "1001",
					"material": "co",
				},
			},
		},
		{
			"UpdateAlarm",
			&mockNATSConnection{},
			&mockAlarmService{
				UpdateOutUpdated: true,
			},
			&mockCameraService{},
			map[string]interface{}{
				"kind": updateAlarm,
				"id":   "aaaa-aaaa",
				"input": map[string]interface{}{
					"siteId":   "1111-1111",
					"serialNo": "1002",
					"material": "smoke",
				},
			},
			map[string]interface{}{
				"kind":    updateAlarm,
				"updated": true,
			},
		},
		{
			"DeleteAlarm",
			&mockNATSConnection{},
			&mockAlarmService{
				DeleteOutDeleted: true,
			},
			&mockCameraService{},
			map[string]interface{}{
				"kind": deleteAlarm,
				"id":   "aaaa-aaaa",
			},
			map[string]interface{}{
				"kind":    deleteAlarm,
				"deleted": true,
			},
		},
		{
			"CreateCamera",
			&mockNATSConnection{},
			&mockAlarmService{},
			&mockCameraService{
				CreateOutCamera: &model.Camera{
					Asset: model.Asset{
						ID:       "bbbb-bbbb",
						SiteID:   "1111-1111",
						SerialNo: "2001",
					},
					Resolution: 921600,
				},
			},
			map[string]interface{}{
				"kind": createCamera,
				"input": map[string]interface{}{
					"siteId":     "1111-1111",
					"serialNo":   "2001",
					"resolution": 921600,
				},
			},
			map[string]interface{}{
				"kind": createCamera,
				"camera": map[string]interface{}{
					"id":         "bbbb-bbbb",
					"siteId":     "1111-1111",
					"serialNo":   "2001",
					"resolution": float64(921600),
				},
			},
		},
		{
			"AllCamera",
			&mockNATSConnection{},
			&mockAlarmService{},
			&mockCameraService{
				AllOutCameras: []model.Camera{
					model.Camera{
						Asset: model.Asset{
							ID:       "bbbb-bbbb",
							SiteID:   "1111-1111",
							SerialNo: "2001",
						},
						Resolution: 921600,
					},
				},
			},
			map[string]interface{}{
				"kind":   allCamera,
				"siteId": "1111-1111",
			},
			map[string]interface{}{
				"kind": allCamera,
				"cameras": []interface{}{
					map[string]interface{}{
						"id":         "bbbb-bbbb",
						"siteId":     "1111-1111",
						"serialNo":   "2001",
						"resolution": float64(921600),
					},
				},
			},
		},
		{
			"GetCamera",
			&mockNATSConnection{},
			&mockAlarmService{},
			&mockCameraService{
				GetOutCamera: &model.Camera{
					Asset: model.Asset{
						ID:       "bbbb-bbbb",
						SiteID:   "1111-1111",
						SerialNo: "2001",
					},
					Resolution: 921600,
				},
			},
			map[string]interface{}{
				"kind": getCamera,
				"id":   "bbbb-bbbb",
			},
			map[string]interface{}{
				"kind": getCamera,
				"camera": map[string]interface{}{
					"id":         "bbbb-bbbb",
					"siteId":     "1111-1111",
					"serialNo":   "2001",
					"resolution": float64(921600),
				},
			},
		},
		{
			"UpdateCamera",
			&mockNATSConnection{},
			&mockAlarmService{},
			&mockCameraService{
				UpdateOutUpdated: true,
			},
			map[string]interface{}{
				"kind": updateCamera,
				"id":   "bbbb-bbbb",
				"input": map[string]interface{}{
					"siteId":     "1111-1111",
					"serialNo":   "2002",
					"resolution": 2073600,
				},
			},
			map[string]interface{}{
				"kind":    updateCamera,
				"updated": true,
			},
		},
		{
			"DeleteCamera",
			&mockNATSConnection{},
			&mockAlarmService{},
			&mockCameraService{
				DeleteOutDeleted: true,
			},
			map[string]interface{}{
				"kind": deleteCamera,
				"id":   "bbbb-bbbb",
			},
			map[string]interface{}{
				"kind":    deleteCamera,
				"deleted": true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()

			nt := &natsTransport{
				logger:        logger,
				metrics:       metrics,
				tracer:        tracer,
				conn:          tc.conn,
				alarmService:  tc.alarmService,
				cameraService: tc.cameraService,
			}

			err := nt.Start()
			assert.Equal(t, tc.conn.QueueSubscribeOutError, err)
			assert.Equal(t, tc.conn.QueueSubscribeOutSubscription, nt.subscription)

			data, err := json.Marshal(tc.request)
			assert.NoError(t, err)

			msg := &nats.Msg{
				Subject: tc.conn.QueueSubscribeInSubject,
				Reply:   "reply_here",
				Data:    data,
			}

			// Simulate the message handling by manually calling the callback function
			callback := tc.conn.QueueSubscribeInCallback
			callback(msg)

			if tc.expectedResponse != nil {
				var response map[string]interface{}
				err = json.Unmarshal(tc.conn.PublishInData, &response)

				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)

				// Verify trace span
				span := tracer.FinishedSpans()[0]
				assert.Equal(t, tc.request["kind"], span.OperationName)
				assert.Equal(t, "NATS", span.Tag("broker"))
				assert.Equal(t, msg.Subject, span.Tag("subject"))
				assert.Equal(t, msg.Reply, span.Tag("reply"))
			}
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name          string
		conn          *mockNATSConnection
		subscription  *nats.Subscription
		expectedError error
	}{
		{
			"Default",
			&mockNATSConnection{},
			&nats.Subscription{},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()

			nt := &natsTransport{
				logger:       logger,
				metrics:      metrics,
				tracer:       tracer,
				conn:         tc.conn,
				subscription: tc.subscription,
			}

			err := nt.Stop(context.Background())
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
