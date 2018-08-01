package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/model"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/opentracing/opentracing-go"
)

const (
	queryGetSwitches = `FOR sw IN switches FILTER sw.siteId == @siteId RETURN sw`
)

// SwitchService implements proto.SwitchServiceServer
type SwitchService struct {
	arango  ArangoService
	logger  *log.Logger
	metrics *metrics.Metrics
	tracer  opentracing.Tracer
}

// NewSwitchService creates a new switch service
func NewSwitchService(arango ArangoService, logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) proto.SwitchServiceServer {
	return &SwitchService{
		arango:  arango,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

// InstallSwitch creates a new switch
func (s *SwitchService) InstallSwitch(ctx context.Context, req *proto.InstallSwitchRequest) (*proto.Switch, error) {
	doc := &model.Switch{
		Key:    uuid.New().String(),
		SiteID: req.GetSiteId(),
		Name:   req.GetName(),
		State:  req.GetState(),
		States: req.GetStates(),
	}

	start := time.Now()
	meta, err := s.arango.CreateDocument(ctx, doc)
	duration := time.Now().Sub(start).Seconds()

	if err != nil {
		s.logger.Error("message", fmt.Sprintf("InstallSwitch failed: %s", err))
		s.metrics.OpLatencyHist.WithLabelValues("install_switch_create_doc", "false").Observe(duration)
		s.metrics.OpLatencySumm.WithLabelValues("install_switch_create_doc", "false").Observe(duration)
		return nil, err
	}

	s.logger.Debug("message", "InstallSwitch succeeded.", "req", req)
	s.metrics.OpLatencyHist.WithLabelValues("install_switch_create_doc", "true").Observe(duration)
	s.metrics.OpLatencySumm.WithLabelValues("install_switch_create_doc", "true").Observe(duration)

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
	key := req.GetId()

	start := time.Now()
	_, err := s.arango.RemoveDocument(ctx, key)
	duration := time.Now().Sub(start).Seconds()

	if err != nil {
		s.logger.Error("message", fmt.Sprintf("RemoveSwitch failed: %s", err))
		s.metrics.OpLatencyHist.WithLabelValues("remove_switch_remove_doc", "false").Observe(duration)
		s.metrics.OpLatencySumm.WithLabelValues("remove_switch_remove_doc", "false").Observe(duration)
		return nil, err
	}

	s.logger.Debug("message", "RemoveSwitch succeeded.", "req", req)
	s.metrics.OpLatencyHist.WithLabelValues("remove_switch_remove_doc", "true").Observe(duration)
	s.metrics.OpLatencySumm.WithLabelValues("remove_switch_remove_doc", "true").Observe(duration)

	return &proto.RemoveSwitchResponse{}, nil
}

// GetSwitch retrieves a switch
func (s *SwitchService) GetSwitch(ctx context.Context, req *proto.GetSwitchRequest) (*proto.Switch, error) {
	key := req.GetId()
	doc := &model.Switch{}

	start := time.Now()
	_, err := s.arango.ReadDocument(ctx, key, doc)
	duration := time.Now().Sub(start).Seconds()

	if err != nil {
		s.logger.Error("message", fmt.Sprintf("GetSwitch failed: %s", err))
		s.metrics.OpLatencyHist.WithLabelValues("get_switch_read_doc", "false").Observe(duration)
		s.metrics.OpLatencySumm.WithLabelValues("get_switch_read_doc", "false").Observe(duration)
		return nil, err
	}

	s.logger.Debug("message", "GetSwitch succeeded.", "req", req)
	s.metrics.OpLatencyHist.WithLabelValues("get_switch_read_doc", "true").Observe(duration)
	s.metrics.OpLatencySumm.WithLabelValues("get_switch_read_doc", "true").Observe(duration)

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
	ctx := stream.Context()
	vars := map[string]interface{}{
		"siteId": req.GetSiteId(),
	}

	start := time.Now()
	cursor, err := s.arango.Query(ctx, queryGetSwitches, vars)
	duration := time.Now().Sub(start).Seconds()

	if err != nil {
		s.logger.Error("message", fmt.Sprintf("GetSwitches failed: %s", err))
		s.metrics.OpLatencyHist.WithLabelValues("get_switches_query", "false").Observe(duration)
		s.metrics.OpLatencySumm.WithLabelValues("get_switches_query", "false").Observe(duration)
		return err
	}
	defer cursor.Close()

	s.metrics.OpLatencyHist.WithLabelValues("get_switches_query", "true").Observe(duration)
	s.metrics.OpLatencySumm.WithLabelValues("get_switches_query", "true").Observe(duration)

	start = time.Now()

	for cursor.HasMore() {
		doc := &model.Switch{}
		_, err := cursor.ReadDocument(ctx, doc)
		if err != nil {
			duration = time.Now().Sub(start).Seconds()
			s.logger.Error("message", fmt.Sprintf("GetSwitches failed: %s", err))
			s.metrics.OpLatencyHist.WithLabelValues("get_switches_send_stream", "false").Observe(duration)
			s.metrics.OpLatencySumm.WithLabelValues("get_switches_send_stream", "false").Observe(duration)
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
			duration = time.Now().Sub(start).Seconds()
			s.logger.Error("message", fmt.Sprintf("GetSwitches failed: %s", err))
			s.metrics.OpLatencyHist.WithLabelValues("get_switches_send_stream", "false").Observe(duration)
			s.metrics.OpLatencySumm.WithLabelValues("get_switches_send_stream", "false").Observe(duration)
			return err
		}
	}

	duration = time.Now().Sub(start).Seconds()
	s.logger.Debug("message", "GetSwitches succeeded.", "req", req)
	s.metrics.OpLatencyHist.WithLabelValues("get_switches_send_stream", "true").Observe(duration)
	s.metrics.OpLatencySumm.WithLabelValues("get_switches_send_stream", "true").Observe(duration)

	return nil
}

// SetSwitch changes the state of a switch
func (s *SwitchService) SetSwitch(ctx context.Context, req *proto.SetSwitchRequest) (*proto.SetSwitchResponse, error) {
	key := req.GetId()
	doc := &model.Switch{
		State: req.GetState(),
	}

	start := time.Now()
	_, err := s.arango.UpdateDocument(ctx, key, doc)
	duration := time.Now().Sub(start).Seconds()

	if err != nil {
		s.logger.Error("message", fmt.Sprintf("SetSwitch failed: %s", err))
		s.metrics.OpLatencyHist.WithLabelValues("set_switch_update_doc", "false").Observe(duration)
		s.metrics.OpLatencySumm.WithLabelValues("set_switch_update_doc", "false").Observe(duration)
		return nil, err
	}

	s.logger.Debug("message", "SetSwitch succeeded.", "req", req)
	s.metrics.OpLatencyHist.WithLabelValues("set_switch_update_doc", "true").Observe(duration)
	s.metrics.OpLatencySumm.WithLabelValues("set_switch_update_doc", "true").Observe(duration)

	return &proto.SetSwitchResponse{}, nil
}
