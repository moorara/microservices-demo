package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/moorara/microservices-demo/services/asset/internal/db"
	"github.com/moorara/microservices-demo/services/asset/internal/model"
	"github.com/moorara/microservices-demo/services/asset/internal/service"
	"github.com/moorara/microservices-demo/services/asset/pkg/log"
	"github.com/moorara/microservices-demo/services/asset/pkg/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

type Entity struct {
	ID   string `json:"id" gorm:"primary_key"`
	Name string
}

func contextWithSpan() context.Context {
	tracer := mocktracer.New()
	span := tracer.StartSpan("test-span")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	return ctx
}

func TestCockroachORM(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

	tests := []struct {
		name   string
		create Entity
		update Entity
	}{
		{
			name:   "Icon",
			create: Entity{Name: "Icon48x48"},
			update: Entity{Name: "Icon64x64"},
		},
		{
			name:   "Font",
			create: Entity{Name: "FontCool"},
			update: Entity{Name: "FontAwesome"},
		},
	}

	logger := log.NewLogger("integration-test", "TestCockroachORM", Config.LogLevel)

	orm, err := db.NewCockroachORM(Config.CockroachAddr, Config.CockroachUser, Config.CockroachPassword, Config.CockroachDatabase, logger)
	assert.NoError(t, err)
	assert.NotNil(t, orm)
	defer orm.Close()

	// Create/migrate database tables
	orm.AutoMigrate(&Entity{})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New().String()
			tc.create.ID = id

			// CREATE
			t.Run("Create", func(t *testing.T) {
				err := orm.Create(&tc.create).Error
				assert.NoError(t, err)
			})

			// UPDATE
			t.Run("Update", func(t *testing.T) {
				err := orm.Model(&tc.update).Where("ID = ?", id).Update(tc.update).Error
				assert.NoError(t, err)
			})

			// READ
			t.Run("Read", func(t *testing.T) {
				var entity Entity
				err := orm.Find(&entity, "ID = ?", id).Error
				assert.NoError(t, err)
				assert.Equal(t, id, entity.ID)
				assert.Equal(t, tc.update.Name, entity.Name)
			})

			// DELETE
			t.Run("Delete", func(t *testing.T) {
				result := orm.Delete(Entity{}, "ID = ?", id)
				assert.NoError(t, result.Error)
				assert.Equal(t, int64(1), result.RowsAffected)
			})
		})
	}
}

func TestAlarmService(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

	tests := []struct {
		name    string
		alarms  []*model.Alarm
		siteIDs []string
		creates []model.AlarmInput
		updates []model.AlarmInput
	}{
		{
			"Simple",
			[]*model.Alarm{},
			[]string{"1111-1111", "2222-2222"},
			[]model.AlarmInput{
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "1001"}, Material: "co"},
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "1002"}, Material: "smoke"},
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "1003"}, Material: "co"},
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "1004"}, Material: "smoke"},
			},
			[]model.AlarmInput{
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "9999"}, Material: "gas"},
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "9998"}, Material: "propane"},
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "9997"}, Material: "gas"},
				model.AlarmInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "9996"}, Material: "propane"},
			},
		},
	}

	logger := log.NewLogger("integration-test", "TestAlarmService", Config.LogLevel)
	metrics := metrics.New("integration-test")
	tracer := mocktracer.New()

	orm, err := db.NewCockroachORM(Config.CockroachAddr, Config.CockroachUser, Config.CockroachPassword, Config.CockroachDatabase, logger)
	assert.NoError(t, err)
	assert.NotNil(t, orm)
	defer orm.Close()

	alarmService := service.NewAlarmService(orm, logger, metrics, tracer)
	assert.NotNil(t, alarmService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("Create", func(t *testing.T) {
				for _, input := range tc.creates {
					ctx := contextWithSpan()
					alarm, err := alarmService.Create(ctx, input)
					assert.NoError(t, err)
					assert.NotEmpty(t, alarm.Asset.ID)
					assert.Equal(t, input.AssetInput.SiteID, alarm.Asset.SiteID)
					assert.Equal(t, input.AssetInput.SerialNo, alarm.Asset.SerialNo)
					assert.Equal(t, input.Material, alarm.Material)
					tc.alarms = append(tc.alarms, alarm)
				}
			})

			t.Run("All", func(t *testing.T) {
				for _, siteID := range tc.siteIDs {
					ctx := contextWithSpan()
					alarms, err := alarmService.All(ctx, siteID)
					assert.NoError(t, err)
					assert.True(t, len(alarms) > 0)
				}
			})

			t.Run("Get", func(t *testing.T) {
				for _, alarm := range tc.alarms {
					ctx := contextWithSpan()
					result, err := alarmService.Get(ctx, alarm.ID)
					assert.NoError(t, err)
					assert.Equal(t, alarm.Asset.ID, result.Asset.ID)
					assert.Equal(t, alarm.Asset.SiteID, result.Asset.SiteID)
					assert.Equal(t, alarm.Asset.SerialNo, result.Asset.SerialNo)
					assert.Equal(t, alarm.Material, result.Material)
				}
			})

			t.Run("Update", func(t *testing.T) {
				for i := range tc.alarms {
					ctx := contextWithSpan()
					result, err := alarmService.Update(ctx, tc.alarms[i].ID, tc.updates[i])
					assert.NoError(t, err)
					assert.True(t, result)
				}
			})

			t.Run("Delete", func(t *testing.T) {
				for _, alarm := range tc.alarms {
					ctx := contextWithSpan()
					result, err := alarmService.Delete(ctx, alarm.ID)
					assert.NoError(t, err)
					assert.True(t, result)
				}
			})
		})
	}
}

func TestCameraService(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

	tests := []struct {
		name    string
		cameras []*model.Camera
		siteIDs []string
		creates []model.CameraInput
		updates []model.CameraInput
	}{
		{
			"Simple",
			[]*model.Camera{},
			[]string{"1111-1111", "2222-2222"},
			[]model.CameraInput{
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "2001"}, Resolution: 1920000},
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "2002"}, Resolution: 1920000},
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "2003"}, Resolution: 1920000},
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "2004"}, Resolution: 1920000},
			},
			[]model.CameraInput{
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "9999"}, Resolution: 4915200},
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "1111-1111", SerialNo: "9998"}, Resolution: 4915200},
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "9997"}, Resolution: 4915200},
				model.CameraInput{AssetInput: model.AssetInput{SiteID: "2222-2222", SerialNo: "9996"}, Resolution: 4915200},
			},
		},
	}

	logger := log.NewLogger("integration-test", "TestCameraService", Config.LogLevel)
	metrics := metrics.New("integration-test")
	tracer := mocktracer.New()

	orm, err := db.NewCockroachORM(Config.CockroachAddr, Config.CockroachUser, Config.CockroachPassword, Config.CockroachDatabase, logger)
	assert.NoError(t, err)
	assert.NotNil(t, orm)
	defer orm.Close()

	cameraService := service.NewCameraService(orm, logger, metrics, tracer)
	assert.NotNil(t, cameraService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("Create", func(t *testing.T) {
				for _, input := range tc.creates {
					ctx := contextWithSpan()
					camera, err := cameraService.Create(ctx, input)
					assert.NoError(t, err)
					assert.NotEmpty(t, camera.Asset.ID)
					assert.Equal(t, input.AssetInput.SiteID, camera.Asset.SiteID)
					assert.Equal(t, input.AssetInput.SerialNo, camera.Asset.SerialNo)
					assert.Equal(t, input.Resolution, camera.Resolution)
					tc.cameras = append(tc.cameras, camera)
				}
			})

			t.Run("All", func(t *testing.T) {
				for _, siteID := range tc.siteIDs {
					ctx := contextWithSpan()
					cameras, err := cameraService.All(ctx, siteID)
					assert.NoError(t, err)
					assert.True(t, len(cameras) > 0)
				}
			})

			t.Run("Get", func(t *testing.T) {
				for _, camera := range tc.cameras {
					ctx := contextWithSpan()
					result, err := cameraService.Get(ctx, camera.ID)
					assert.NoError(t, err)
					assert.Equal(t, camera.Asset.ID, result.Asset.ID)
					assert.Equal(t, camera.Asset.SiteID, result.Asset.SiteID)
					assert.Equal(t, camera.Asset.SerialNo, result.Asset.SerialNo)
					assert.Equal(t, camera.Resolution, result.Resolution)
				}
			})

			t.Run("Update", func(t *testing.T) {
				for i := range tc.cameras {
					ctx := contextWithSpan()
					result, err := cameraService.Update(ctx, tc.cameras[i].ID, tc.updates[i])
					assert.NoError(t, err)
					assert.True(t, result)
				}
			})

			t.Run("Delete", func(t *testing.T) {
				for _, camera := range tc.cameras {
					ctx := contextWithSpan()
					result, err := cameraService.Delete(ctx, camera.ID)
					assert.NoError(t, err)
					assert.True(t, result)
				}
			})
		})
	}
}
