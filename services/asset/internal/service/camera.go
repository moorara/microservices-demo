package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/moorara/microservices-demo/services/asset/internal/db"
	"github.com/moorara/microservices-demo/services/asset/internal/model"
	"github.com/moorara/microservices-demo/services/asset/pkg/log"
	"github.com/moorara/microservices-demo/services/asset/pkg/metrics"
	"github.com/opentracing/opentracing-go"

	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
)

type (
	// CameraService is the service for Camera model CRUD
	CameraService interface {
		Create(ctx context.Context, input model.CameraInput) (*model.Camera, error)
		All(ctx context.Context, siteID string) ([]model.Camera, error)
		Get(ctx context.Context, id string) (*model.Camera, error)
		Update(ctx context.Context, id string, input model.CameraInput) (bool, error)
		Delete(ctx context.Context, id string) (bool, error)
	}

	cameraService struct {
		orm     db.ORM
		logger  *log.Logger
		metrics *metrics.Metrics
		tracer  opentracing.Tracer
	}
)

// NewCameraService creates a new CameraService object
func NewCameraService(orm db.ORM, logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) CameraService {
	// Migrate the database table schema
	orm.AutoMigrate(model.Camera{})

	return &cameraService{
		orm:     orm,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (s *cameraService) exec(ctx context.Context, op, query string, fn func() error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := s.tracer.StartSpan(op, opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	ext.DBType.Set(span, "sql")
	ext.DBStatement.Set(span, query)
	span.LogFields(opentracingLog.String("event", op))

	start := time.Now()
	err := fn()
	latency := time.Now().Sub(start).Seconds()

	success := "true"
	if err != nil {
		success = "false"
		s.logger.Error("message", fmt.Sprintf("%s failed: %s", op, err))
		span.LogFields(opentracingLog.String("message", err.Error()))
	} else {
		s.logger.Debug("message", fmt.Sprintf("%s succeeded.", op))
		span.LogFields(opentracingLog.String("message", "successful!"))
	}

	s.metrics.OpLatencyHist.WithLabelValues(op, success).Observe(latency)
	s.metrics.OpLatencySumm.WithLabelValues(op, success).Observe(latency)
}

func (s *cameraService) Create(ctx context.Context, input model.CameraInput) (*model.Camera, error) {
	var err error

	camera := &model.Camera{
		Asset: model.Asset{
			ID:       uuid.New().String(),
			SiteID:   input.SiteID,
			SerialNo: input.SerialNo,
		},
		Resolution: input.Resolution,
	}

	s.exec(ctx, "create_camera", "gorm.Create", func() error {
		err = s.orm.Create(&camera).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return camera, nil
}

func (s *cameraService) All(ctx context.Context, siteID string) ([]model.Camera, error) {
	var err error
	var cameras []model.Camera

	s.exec(ctx, "all_cameras", "gorm.Find", func() error {
		err = s.orm.Find(&cameras, "site_id = ?", siteID).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return cameras, nil
}

func (s *cameraService) Get(ctx context.Context, id string) (*model.Camera, error) {
	var err error
	camera := &model.Camera{}

	s.exec(ctx, "get_camera", "gorm.Find", func() error {
		err = s.orm.Find(camera, "id = ?", id).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return camera, nil
}

func (s *cameraService) Update(ctx context.Context, id string, input model.CameraInput) (bool, error) {
	var result *gorm.DB

	camera := &model.Camera{
		Asset: model.Asset{
			ID:       id,
			SiteID:   input.SiteID,
			SerialNo: input.SerialNo,
		},
		Resolution: input.Resolution,
	}

	s.exec(ctx, "update_camera", "gorm.Model.Where.Update", func() error {
		result = s.orm.Model(camera).Where("id = ?", id).Update(camera)
		return result.Error
	})

	if err := result.Error; err != nil {
		return false, err
	}

	return result.RowsAffected == 1, nil
}

func (s *cameraService) Delete(ctx context.Context, id string) (bool, error) {
	var result *gorm.DB

	s.exec(ctx, "delete_camera", "gorm.Delete", func() error {
		result = s.orm.Delete(model.Camera{}, "id = ?", id)
		return result.Error
	})

	if err := result.Error; err != nil {
		return false, err
	}

	return result.RowsAffected == 1, nil
}
