package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/moorara/microservices-demo/services/asset-service/internal/db"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

const (
	CockroachDatabase = "assets"
)

type Asset struct {
	ID   string `json:"id" gorm:"primary_key"`
	Name string
}

func TestCockroachDB(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

	tests := []struct {
		name   string
		create Asset
		update Asset
	}{
		{
			name:   "Icon",
			create: Asset{Name: "Icon48x48"},
			update: Asset{Name: "Icon64x64"},
		},
		{
			name:   "Font",
			create: Asset{Name: "FontCool"},
			update: Asset{Name: "FontAwesome"},
		},
	}

	logger := log.NewLogger("integration-test", "singleton", Config.LogLevel)
	orm := db.NewCockroachORM(Config.CockroachAddr, Config.CockroachUser, Config.CockroachPassword, CockroachDatabase, logger)
	assert.NotNil(t, orm)
	defer orm.Close()

	// Create/migrate database tables
	orm.AutoMigrate(&Asset{})

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
				var asset Asset
				err := orm.Find(&asset, "ID = ?", id).Error
				assert.NoError(t, err)
				assert.Equal(t, id, asset.ID)
				assert.Equal(t, tc.update.Name, asset.Name)
			})

			// DELETE
			t.Run("Delete", func(t *testing.T) {
				result := orm.Delete(Asset{}, "ID = ?", id)
				assert.NoError(t, result.Error)
				assert.Equal(t, int64(1), result.RowsAffected)
			})
		})
	}
}
