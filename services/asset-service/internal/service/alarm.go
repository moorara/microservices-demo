package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/moorara/microservices-demo/services/asset-service/internal/db"
	"github.com/moorara/microservices-demo/services/asset-service/internal/model"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	opentracingLog "github.com/opentracing/opentracing-go/log"
)

type (
	// AlarmService is the service for Alarm model CRUD
	AlarmService interface {
		Create(ctx context.Context, input model.AlarmInput) (*model.Alarm, error)
		All(ctx context.Context, siteID string) ([]model.Alarm, error)
		Get(ctx context.Context, id string) (*model.Alarm, error)
		Update(ctx context.Context, id string, input model.AlarmInput) (bool, error)
		Delete(ctx context.Context, id string) (bool, error)
	}

	alarmService struct {
		orm     db.ORM
		logger  *log.Logger
		metrics *metrics.Metrics
		tracer  opentracing.Tracer
	}
)

// NewAlarmService creates a new AlarmService object
func NewAlarmService(orm db.ORM, logger *log.Logger, metrics *metrics.Metrics, tracer opentracing.Tracer) AlarmService {
	// Migrate the database table schema
	orm.AutoMigrate(model.Alarm{})

	return &alarmService{
		orm:     orm,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (s *alarmService) exec(ctx context.Context, op, query string, fn func() error) {
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

func (s *alarmService) Create(ctx context.Context, input model.AlarmInput) (*model.Alarm, error) {
	var err error

	alarm := &model.Alarm{
		Asset: model.Asset{
			ID:       uuid.New().String(),
			SiteID:   input.SiteID,
			SerialNo: input.SerialNo,
		},
		Material: input.Material,
	}

	s.exec(ctx, "create_alarm", "gorm.Create", func() error {
		err = s.orm.Create(&alarm).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return alarm, nil
}

func (s *alarmService) All(ctx context.Context, siteID string) ([]model.Alarm, error) {
	var err error
	var alarms []model.Alarm

	s.exec(ctx, "all_alarms", "gorm.Find", func() error {
		err = s.orm.Find(&alarms, "site_id = ?", siteID).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return alarms, nil
}

func (s *alarmService) Get(ctx context.Context, id string) (*model.Alarm, error) {
	var err error
	alarm := &model.Alarm{}

	s.exec(ctx, "get_alarm", "gorm.Find", func() error {
		err = s.orm.Find(alarm, "id = ?", id).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return alarm, nil
}

func (s *alarmService) Update(ctx context.Context, id string, input model.AlarmInput) (bool, error) {
	var result *gorm.DB

	alarm := &model.Alarm{
		Asset: model.Asset{
			ID:       id,
			SiteID:   input.SiteID,
			SerialNo: input.SerialNo,
		},
		Material: input.Material,
	}

	s.exec(ctx, "update_alarm", "gorm.Model.Where.Update", func() error {
		result = s.orm.Model(alarm).Where("id = ?", id).Update(alarm)
		return result.Error
	})

	if err := result.Error; err != nil {
		return false, err
	}

	return result.RowsAffected == 1, nil
}

func (s *alarmService) Delete(ctx context.Context, id string) (bool, error) {
	var result *gorm.DB

	s.exec(ctx, "delete_alarm", "gorm.Delete", func() error {
		result = s.orm.Delete(model.Alarm{}, "id = ?", id)
		return result.Error
	})

	if err := result.Error; err != nil {
		return false, err
	}

	return result.RowsAffected == 1, nil
}
