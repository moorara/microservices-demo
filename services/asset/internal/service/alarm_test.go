package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/moorara/microservices-demo/services/asset/internal/db"
	"github.com/moorara/microservices-demo/services/asset/internal/model"
	"github.com/moorara/microservices-demo/services/asset/pkg/log"
	"github.com/moorara/microservices-demo/services/asset/pkg/metrics"
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

			// Verify trace span
			span := tracer.FinishedSpans()[0]
			assert.Equal(t, "create_alarm", span.OperationName)
			assert.Equal(t, "sql", span.Tag("db.type"))
			assert.Equal(t, "gorm.Create", span.Tag("db.statement"))
			assert.Equal(t, "event", span.Logs()[0].Fields[0].Key)
			assert.Equal(t, "create_alarm", span.Logs()[0].Fields[0].ValueString)
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

			// Verify trace span
			span := tracer.FinishedSpans()[0]
			assert.Equal(t, "all_alarms", span.OperationName)
			assert.Equal(t, "sql", span.Tag("db.type"))
			assert.Equal(t, "gorm.Find", span.Tag("db.statement"))
			assert.Equal(t, "event", span.Logs()[0].Fields[0].Key)
			assert.Equal(t, "all_alarms", span.Logs()[0].Fields[0].ValueString)
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

			// Verify trace span
			span := tracer.FinishedSpans()[0]
			assert.Equal(t, "get_alarm", span.OperationName)
			assert.Equal(t, "sql", span.Tag("db.type"))
			assert.Equal(t, "gorm.Find", span.Tag("db.statement"))
			assert.Equal(t, "event", span.Logs()[0].Fields[0].Key)
			assert.Equal(t, "get_alarm", span.Logs()[0].Fields[0].ValueString)
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

			// Verify trace span
			span := tracer.FinishedSpans()[0]
			assert.Equal(t, "update_alarm", span.OperationName)
			assert.Equal(t, "sql", span.Tag("db.type"))
			assert.Equal(t, "gorm.Model.Where.Update", span.Tag("db.statement"))
			assert.Equal(t, "event", span.Logs()[0].Fields[0].Key)
			assert.Equal(t, "update_alarm", span.Logs()[0].Fields[0].ValueString)
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

			// Verify trace span
			span := tracer.FinishedSpans()[0]
			assert.Equal(t, "delete_alarm", span.OperationName)
			assert.Equal(t, "sql", span.Tag("db.type"))
			assert.Equal(t, "gorm.Delete", span.Tag("db.statement"))
			assert.Equal(t, "event", span.Logs()[0].Fields[0].Key)
			assert.Equal(t, "delete_alarm", span.Logs()[0].Fields[0].ValueString)
		})
	}
}
