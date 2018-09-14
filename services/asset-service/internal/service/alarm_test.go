package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/moorara/microservices-demo/services/asset-service/internal/db"
	"github.com/moorara/microservices-demo/services/asset-service/internal/model"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/metrics"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestNewAlarmService(t *testing.T) {
	tests := []struct {
		name string
		orm  db.ORM
	}{
		{
			"Default",
			&mockORM{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := NewAlarmService(tc.orm, logger, metrics, tracer)

			assert.NotNil(t, service)
		})
	}
}

func TestAlarmServiceCreate(t *testing.T) {
	tests := []struct {
		name          string
		orm           db.ORM
		ctx           context.Context
		input         model.AlarmInput
		expectedError error
	}{
		{
			"DatabaseError",
			&mockORM{
				CreateOutDB: &gorm.DB{
					Error: errors.New("create error"),
				},
			},
			contextWithSpan(),
			model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "1001"}, Material: "smoke"},
			errors.New("create error"),
		},
		{
			"Success",
			&mockORM{
				CreateOutDB: &gorm.DB{},
			},
			contextWithSpan(),
			model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "1001"}, Material: "smoke"},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &alarmService{tc.orm, logger, metrics, tracer}

			_, err := service.Create(tc.ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAlarmServiceAll(t *testing.T) {
	tests := []struct {
		name          string
		orm           db.ORM
		ctx           context.Context
		siteID        string
		expectedError error
	}{
		{
			"DatabaseError",
			&mockORM{
				FindOutDB: &gorm.DB{
					Error: errors.New("find error"),
				},
			},
			contextWithSpan(),
			"1111-1111",
			errors.New("find error"),
		},
		{
			"Success",
			&mockORM{
				FindOutDB: &gorm.DB{},
			},
			contextWithSpan(),
			"1111-1111",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &alarmService{tc.orm, logger, metrics, tracer}

			_, err := service.All(tc.ctx, tc.siteID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAlarmServiceGet(t *testing.T) {
	tests := []struct {
		name          string
		orm           db.ORM
		ctx           context.Context
		id            string
		expectedError error
	}{
		{
			"DatabaseError",
			&mockORM{
				FindOutDB: &gorm.DB{
					Error: errors.New("find error"),
				},
			},
			contextWithSpan(),
			"aaaa-aaaa",
			errors.New("find error"),
		},
		{
			"Success",
			&mockORM{
				FindOutDB: &gorm.DB{},
			},
			contextWithSpan(),
			"aaaa-aaaa",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &alarmService{tc.orm, logger, metrics, tracer}

			_, err := service.Get(tc.ctx, tc.id)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAlarmServiceUpdate(t *testing.T) {
	tests := []struct {
		name           string
		orm            db.ORM
		ctx            context.Context
		id             string
		input          model.AlarmInput
		expectedError  error
		expectedResult bool
	}{
		/* {
			"DatabaseError",
			&mockORM{},
			CreateContextWithSpan(),
			"aaaa-aaaa",
			model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "1002"}, Material: "co"},
			errors.New(""),
			false,
		},
		{
			"Success",
			&mockORM{},
			CreateContextWithSpan(),
			"aaaa-aaaa",
			model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "1002"}, Material: "co"},
			nil,
			true,
		}, */
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &alarmService{tc.orm, logger, metrics, tracer}

			result, err := service.Update(tc.ctx, tc.id, tc.input)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAlarmServiceDelete(t *testing.T) {
	tests := []struct {
		name           string
		orm            db.ORM
		ctx            context.Context
		id             string
		expectedError  error
		expectedResult bool
	}{
		{
			"DatabaseError",
			&mockORM{
				DeleteOutDB: &gorm.DB{
					Error: errors.New("delete error"),
				},
			},
			contextWithSpan(),
			"aaaa-aaaa",
			errors.New("delete error"),
			false,
		},
		{
			"Success",
			&mockORM{
				DeleteOutDB: &gorm.DB{
					RowsAffected: 1,
				},
			},
			contextWithSpan(),
			"aaaa-aaaa",
			nil,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &alarmService{tc.orm, logger, metrics, tracer}

			result, err := service.Delete(tc.ctx, tc.id)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
