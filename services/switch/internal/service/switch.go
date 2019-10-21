package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/moorara/microservices-demo/services/switch/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch/internal/model"
	"github.com/moorara/microservices-demo/services/switch/internal/proto"
	"github.com/moorara/microservices-demo/services/switch/pkg/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"

	arango "github.com/arangodb/go-driver"
	opentracingLog "github.com/opentracing/opentracing-go/log"
)

const queryGetSwitches = `FOR sw IN switches FILTER sw.siteId == @siteId RETURN sw`

type (
	callback func() error

	// SwitchService implements proto.SwitchServiceServer
	SwitchService struct {
		arango  ArangoService
		logger  *log.Logger
		metrics *metrics.Metrics
		tracer  opentracing.Tracer
	}
)

// NewSwitchService creates a new switch service
func NewSwitchService(arango ArangoService, logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) proto.SwitchServiceServer {
	return &SwitchService{
		arango:  arango,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (s *SwitchService) extractParentSpanContext(ctx context.Context) (opentracing.SpanContext, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := meta.Get("span.context")
		if len(vals) == 1 {
			data := make(map[string]string)
			err := json.Unmarshal([]byte(vals[0]), &data)
			if err != nil {
				return nil, err
			}

			carrier := opentracing.TextMapCarrier(data)
			return s.tracer.Extract(opentracing.TextMap, carrier)
		}
	}

	return nil, nil
}

func (s *SwitchService) exec(ctx context.Context, req interface{}, op, query string, fn callback) {
	var span opentracing.Span
	parentSpanContext, _ := s.extractParentSpanContext(ctx)
	if parentSpanContext == nil {
		span = s.tracer.StartSpan(op)
	} else {
		span = s.tracer.StartSpan(op, opentracing.ChildOf(parentSpanContext))
	}

	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	defer span.Finish()
	ext.DBType.Set(span, "arango")
	ext.DBStatement.Set(span, query)
	span.LogFields(opentracingLog.String("event", op))

	start := time.Now()
	err := fn()
	latency := time.Since(start).Seconds()

	success := "true"
	if err != nil {
		success = "false"
		s.logger.Error("message", fmt.Sprintf("%s failed: %s", op, err))
		span.LogFields(opentracingLog.String("message", err.Error()))
	} else {
		s.logger.Debug("message", fmt.Sprintf("%s succeeded.", op), "req", req)
		span.LogFields(opentracingLog.String("message", "successful!"))
	}

	s.metrics.OpLatencyHist.WithLabelValues(op, success).Observe(latency)
	s.metrics.OpLatencySumm.WithLabelValues(op, success).Observe(latency)
}

// InstallSwitch creates a new switch
func (s *SwitchService) InstallSwitch(ctx context.Context, req *proto.InstallSwitchRequest) (*proto.Switch, error) {
	var err error
	var meta arango.DocumentMeta

	doc := &model.Switch{
		Key:    uuid.New().String(),
		SiteID: req.GetSiteId(),
		Name:   req.GetName(),
		State:  req.GetState(),
		States: req.GetStates(),
	}

	s.exec(ctx, req, "InstallSwitch_CreateDocument", "CreateDocument", func() error {
		meta, err = s.arango.CreateDocument(ctx, doc)
		return err
	})

	if err != nil {
		return nil, err
	}

	doc.ID = meta.ID.String()
	doc.Key = meta.Key
	doc.Rev = meta.Rev

	return &proto.Switch{
		Id:     meta.Key,
		SiteId: doc.SiteID,
		Name:   doc.Name,
		State:  doc.State,
		States: doc.States,
	}, nil
}

// RemoveSwitch deletes a switch
func (s *SwitchService) RemoveSwitch(ctx context.Context, req *proto.RemoveSwitchRequest) (*proto.RemoveSwitchResponse, error) {
	var err error
	key := req.GetId()

	s.exec(ctx, req, "RemoveSwitch_RemoveDocument", "RemoveDocument", func() error {
		_, err = s.arango.RemoveDocument(ctx, key)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &proto.RemoveSwitchResponse{}, nil
}

// GetSwitch retrieves a switch
func (s *SwitchService) GetSwitch(ctx context.Context, req *proto.GetSwitchRequest) (*proto.Switch, error) {
	var err error
	key := req.GetId()
	doc := &model.Switch{}

	s.exec(ctx, req, "GetSwitch_ReadDocument", "ReadDocument", func() error {
		_, err = s.arango.ReadDocument(ctx, key, doc)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &proto.Switch{
		Id:     doc.Key,
		SiteId: doc.SiteID,
		Name:   doc.Name,
		State:  doc.State,
		States: doc.States,
	}, nil
}

// GetSwitches retrieves a group of switches
func (s *SwitchService) GetSwitches(req *proto.GetSwitchesRequest, stream proto.SwitchService_GetSwitchesServer) error {
	var err error
	var cursor arango.Cursor

	ctx := stream.Context()
	vars := map[string]interface{}{
		"siteId": req.GetSiteId(),
	}

	s.exec(ctx, req, "GetSwitches_Query", queryGetSwitches, func() error {
		cursor, err = s.arango.Query(ctx, queryGetSwitches, vars)
		return err
	})

	if err != nil {
		return err
	}

	defer cursor.Close()

	s.exec(ctx, req, "GetSwitches_ReadDocument_Send", "ReadDocument", func() error {
		for cursor.HasMore() {
			doc := &model.Switch{}
			_, err = cursor.ReadDocument(ctx, doc)
			if err != nil {
				return err
			}

			sw := &proto.Switch{
				Id:     doc.Key,
				SiteId: doc.SiteID,
				Name:   doc.Name,
				State:  doc.State,
				States: doc.States,
			}

			err = stream.Send(sw)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// SetSwitch changes the state of a switch
func (s *SwitchService) SetSwitch(ctx context.Context, req *proto.SetSwitchRequest) (*proto.SetSwitchResponse, error) {
	var err error
	key := req.GetId()
	doc := &model.Switch{
		State: req.GetState(),
	}

	s.exec(ctx, req, "SetSwitch_UpdateDocument", "UpdateDocument", func() error {
		_, err = s.arango.UpdateDocument(ctx, key, doc)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &proto.SetSwitchResponse{}, nil
}
