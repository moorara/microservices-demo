package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/moorara/microservices-demo/services/switch-service/internal/metrics"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/log"
	"github.com/opentracing/opentracing-go"
)

// SwitchService implements proto.SwitchServiceServer
type SwitchService struct {
	logger  *log.Logger
	metrics *metrics.Metrics
	tracer  opentracing.Tracer
}

// NewSwitchService creates a new switch service
func NewSwitchService(logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) proto.SwitchServiceServer {
	return &SwitchService{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

// InstallSwitch creates a new switch
func (s *SwitchService) InstallSwitch(ctx context.Context, req *proto.InstallSwitchRequest) (*proto.Switch, error) {
	s.logger.Info("message", "InstallSwitch", "req", req)

	return &proto.Switch{
		Id:     uuid.New().String(),
		SiteId: req.SiteId,
		Name:   req.Name,
		State:  req.State,
		States: req.States,
	}, nil
}

// RemoveSwitch deletes a switch
func (s *SwitchService) RemoveSwitch(ctx context.Context, req *proto.RemoveSwitchRequest) (*proto.RemoveSwitchResponse, error) {
	s.logger.Info("message", "RemoveSwitch", "req", req)
	return &proto.RemoveSwitchResponse{}, nil
}

// GetSwitch retrieves a switch
func (s *SwitchService) GetSwitch(ctx context.Context, req *proto.GetSwitchRequest) (*proto.Switch, error) {
	s.logger.Info("message", "GetSwitch", "req", req)
	return &proto.Switch{Id: req.Id}, nil
}

// GetSwitches retrieves a group of switches
func (s *SwitchService) GetSwitches(req *proto.GetSwitchesRequest, stream proto.SwitchService_GetSwitchesServer) error {
	s.logger.Info("message", "GetSwitches", "req", req)
	stream.Send(&proto.Switch{SiteId: req.SiteId})
	return nil
}

// SetSwitch changes the state of a switch
func (s *SwitchService) SetSwitch(ctx context.Context, req *proto.SetSwitchRequest) (*proto.SetSwitchResponse, error) {
	s.logger.Info("message", "SetSwitch", "req", req)
	return &proto.SetSwitchResponse{}, nil
}
