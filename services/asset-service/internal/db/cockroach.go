package db

import (
	"fmt"

	// Required for initialization
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"github.com/moorara/microservices-demo/services/asset-service/pkg/log"
)

type (
	// ORM is the interface for gorm.DB
	ORM interface {
		AutoMigrate(values ...interface{}) *gorm.DB
		Close() error
		Create(value interface{}) *gorm.DB
		Delete(value interface{}, where ...interface{}) *gorm.DB
		Find(out interface{}, where ...interface{}) *gorm.DB
		LogMode(enable bool) *gorm.DB
		Model(value interface{}) *gorm.DB
		Preload(column string, conditions ...interface{}) *gorm.DB
		Update(attrs ...interface{}) *gorm.DB
		Where(query interface{}, args ...interface{}) *gorm.DB
	}

	gormLogger struct {
		logger *log.Logger
	}
)

func (l *gormLogger) Print(values ...interface{}) {
	l.logger.Debug("message", values[1])
}

// NewCockroachORM creates a new DB for CockroachDB
func NewCockroachORM(addr, user, password, database string, logger *log.Logger) (ORM, error) {
	var url string
	if user != "" && password != "" {
		url = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, addr, database)
	} else if user != "" {
		url = fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", user, addr, database)
	} else {
		url = fmt.Sprintf("postgres://%s/%s?sslmode=disable", addr, database)
	}

	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.LogMode(false)
	db.SetLogger(&gormLogger{logger})

	return db, nil
}
