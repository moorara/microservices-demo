package service

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func contextWithSpan() context.Context {
	tracer := mocktracer.New()
	span := tracer.StartSpan("mock-span")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	return ctx
}

type mockORM struct {
	AutoMigrateCalled   bool
	AutoMigrateInValues []interface{}
	AutoMigrateOutDB    *gorm.DB

	CloseCalled   bool
	CloseOutError error

	CreateCalled  bool
	CreateInValue interface{}
	CreateOutDB   *gorm.DB

	DeleteCalled  bool
	DeleteInValue interface{}
	DeleteInWhere []interface{}
	DeleteOutDB   *gorm.DB

	FindCalled  bool
	FindInOut   interface{}
	FindInWhere []interface{}
	FindOutDB   *gorm.DB

	LogModeCalled   bool
	LogModeInEnable bool
	LogModeOutDB    *gorm.DB

	ModelCalled  bool
	ModelInValue interface{}
	ModelOutDB   *gorm.DB

	PreloadCalled       bool
	PreloadInColumn     string
	PreloadInConditions []interface{}
	PreloadOutDB        *gorm.DB

	UpdateCalled  bool
	UpdateInAttrs []interface{}
	UpdateOutDB   *gorm.DB

	WhereCalled  bool
	WhereInQuery interface{}
	WhereInArgs  []interface{}
	WhereOutDB   *gorm.DB
}

func (m *mockORM) AutoMigrate(values ...interface{}) *gorm.DB {
	m.AutoMigrateCalled = true
	m.AutoMigrateInValues = values
	return m.AutoMigrateOutDB
}

func (m *mockORM) Close() error {
	m.CloseCalled = true
	return m.CloseOutError
}

func (m *mockORM) Create(value interface{}) *gorm.DB {
	m.CreateCalled = true
	m.CreateInValue = value
	return m.CreateOutDB
}

func (m *mockORM) Delete(value interface{}, where ...interface{}) *gorm.DB {
	m.DeleteCalled = true
	m.DeleteInValue = value
	m.DeleteInWhere = where
	return m.DeleteOutDB
}

func (m *mockORM) Find(out interface{}, where ...interface{}) *gorm.DB {
	m.FindCalled = true
	m.FindInOut = out
	m.FindInWhere = where
	return m.FindOutDB
}

func (m *mockORM) LogMode(enable bool) *gorm.DB {
	m.LogModeCalled = true
	m.LogModeInEnable = enable
	return m.LogModeOutDB
}

func (m *mockORM) Model(value interface{}) *gorm.DB {
	m.ModelCalled = true
	m.ModelInValue = value
	return m.ModelOutDB
}

func (m *mockORM) Preload(column string, conditions ...interface{}) *gorm.DB {
	m.PreloadCalled = true
	m.PreloadInColumn = column
	m.PreloadInConditions = conditions
	return m.PreloadOutDB
}

func (m *mockORM) Update(attrs ...interface{}) *gorm.DB {
	m.UpdateCalled = true
	m.UpdateInAttrs = attrs
	return m.UpdateOutDB
}

func (m *mockORM) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.WhereCalled = true
	m.WhereInQuery = query
	m.WhereInArgs = args
	return m.WhereOutDB
}
