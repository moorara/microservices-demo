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

func TestNewCameraService(t *testing.T) {
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
			service := NewCameraService(tc.orm, logger, metrics, tracer)

			assert.NotNil(t, service)
		})
	}
}

func TestCameraServiceCreate(t *testing.T) {
	tests := []struct {
		name          string
		orm           db.ORM
		ctx           context.Context
		input         model.CameraInput
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
			model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "2001"}, Resolution: 1920000},
			errors.New("create error"),
		},
		{
			"Success",
			&mockORM{
				CreateOutDB: &gorm.DB{},
			},
			contextWithSpan(),
			model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "2001"}, Resolution: 1920000},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &cameraService{tc.orm, logger, metrics, tracer}

			_, err := service.Create(tc.ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCameraServiceAll(t *testing.T) {
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
			service := &cameraService{tc.orm, logger, metrics, tracer}

			_, err := service.All(tc.ctx, tc.siteID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCameraServiceGet(t *testing.T) {
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
			"bbbb-bbbb",
			errors.New("find error"),
		},
		{
			"Success",
			&mockORM{
				FindOutDB: &gorm.DB{},
			},
			contextWithSpan(),
			"bbbb-bbbb",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &cameraService{tc.orm, logger, metrics, tracer}

			_, err := service.Get(tc.ctx, tc.id)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCameraServiceUpdate(t *testing.T) {
	tests := []struct {
		name           string
		orm            db.ORM
		ctx            context.Context
		id             string
		input          model.CameraInput
		expectedError  error
		expectedResult bool
	}{
		/* {
			"DatabaseError",
			&mockORM{},
			CreateContextWithSpan(),
			"bbbb-bbbb",
			model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "2002"}, Resolution: 4915200},
			errors.New(""),
			false,
		},
		{
			"Success",
			&mockORM{},
			CreateContextWithSpan(),
			"bbbb-bbbb",
			model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "2002"}, Resolution: 4915200},
			nil,
			true,
		}, */
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &cameraService{tc.orm, logger, metrics, tracer}

			result, err := service.Update(tc.ctx, tc.id, tc.input)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCameraServiceDelete(t *testing.T) {
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
			"bbbb-bbbb",
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
			"bbbb-bbbb",
			nil,
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			metrics := metrics.New("unit-test")
			tracer := mocktracer.New()
			service := &cameraService{tc.orm, logger, metrics, tracer}

			result, err := service.Delete(tc.ctx, tc.id)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
