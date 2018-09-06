package service

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func CreateContextWithSpan() context.Context {
	tracer := mocktracer.New()
	span := tracer.StartSpan("mock-span")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	return ctx
}

func TestNewSensorManager(t *testing.T) {
	tests := []struct {
		name        string
		postgresURL string
	}{
		{
			"WithoutUserPass",
			"postgres://localhost",
		},
		{
			"WithUserPass",
			"postgres://root:pass@localhost",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := NewPostgresDB(tc.postgresURL)
			logger := log.NewNopLogger()
			tracer := mocktracer.New()
			m := NewSensorManager(db, logger, tracer)

			assert.NotNil(t, m)
		})
	}
}

func TestSensorManagerCreate(t *testing.T) {
	tests := []struct {
		name          string
		dbError       error
		dbResult      driver.Result
		context       context.Context
		sensorSiteID  string
		sensorName    string
		sensorUnit    string
		sensorMinSafe float64
		sensorMaxSafe float64
		expectError   bool
	}{
		{
			"DatabaseError",
			errors.New("db error"),
			nil,
			CreateContextWithSpan(),
			"", "", "", 0.0, 0.0,
			true,
		},
		{
			"CreateSensor",
			nil,
			sqlmock.NewResult(0, 1),
			CreateContextWithSpan(),
			"1111-aaaa", "temperature", "celsius", -30.0, 30.0,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
				tracer: mocktracer.New(),
			}

			// Mock SQL query
			expect := mock.ExpectExec(`INSERT INTO sensors`)
			if tc.dbError != nil {
				expect.WillReturnError(tc.dbError)
			} else {
				expect.WillReturnResult(tc.dbResult)
			}

			sensor, err := m.Create(tc.context, tc.sensorSiteID, tc.sensorName, tc.sensorUnit, tc.sensorMinSafe, tc.sensorMaxSafe)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, sensor)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, sensor.ID)
				assert.Equal(t, tc.sensorSiteID, sensor.SiteID)
				assert.Equal(t, tc.sensorName, sensor.Name)
				assert.Equal(t, tc.sensorUnit, sensor.Unit)
				assert.Equal(t, tc.sensorMinSafe, sensor.MinSafe)
				assert.Equal(t, tc.sensorMaxSafe, sensor.MaxSafe)
			}
		})
	}
}

func TestSensorManagerAll(t *testing.T) {
	tests := []struct {
		name            string
		dbError         error
		dbRows          [][]driver.Value
		context         context.Context
		siteID          string
		expectError     bool
		expectedSensors []Sensor
	}{
		{
			"DatabaseError",
			errors.New("db error"),
			nil,
			CreateContextWithSpan(),
			"",
			true,
			nil,
		},
		{
			"AllSensors",
			nil,
			[][]driver.Value{
				[]driver.Value{"2222-bbbb", "1111-aaaa", "temperature", "fahrenheit", -22.0, 86.0},
				[]driver.Value{"3333-cccc", "1111-aaaa", "pressure", "pascal", 50000.0, 100000.0},
			},
			CreateContextWithSpan(),
			"1111-aaaa",
			false,
			[]Sensor{
				Sensor{"2222-bbbb", "1111-aaaa", "temperature", "fahrenheit", -22.0, 86.0},
				Sensor{"3333-cccc", "1111-aaaa", "pressure", "pascal", 50000.0, 100000.0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
				tracer: mocktracer.New(),
			}

			// Mock SQL query
			expect := mock.ExpectQuery(`SELECT id, site_id, name, unit, min_safe, max_safe FROM sensors`)
			if tc.dbError != nil {
				expect.WillReturnError(tc.dbError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "site_id", "name", "unit", "min_safe", "max_safe"})
				for _, row := range tc.dbRows {
					rows.AddRow(row...)
				}
				expect.WillReturnRows(rows)
			}

			sensors, err := m.All(tc.context, tc.siteID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, sensors)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sensors)

				for i, sensor := range sensors {
					assert.Equal(t, tc.expectedSensors[i].ID, sensor.ID)
					assert.Equal(t, tc.expectedSensors[i].SiteID, sensor.SiteID)
					assert.Equal(t, tc.expectedSensors[i].Name, sensor.Name)
					assert.Equal(t, tc.expectedSensors[i].Unit, sensor.Unit)
					assert.Equal(t, tc.expectedSensors[i].MinSafe, sensor.MinSafe)
					assert.Equal(t, tc.expectedSensors[i].MaxSafe, sensor.MaxSafe)
				}
			}
		})
	}
}

func TestSensorManagerGet(t *testing.T) {
	tests := []struct {
		name           string
		dbError        error
		dbRow          []driver.Value
		context        context.Context
		id             string
		expectError    bool
		expectedSensor *Sensor
	}{
		{
			"DatabaseError",
			errors.New("db error"),
			nil,
			CreateContextWithSpan(),
			"",
			true,
			nil,
		},
		{
			"GetSensor",
			nil,
			[]driver.Value{"2222-bbbb", "1111-aaaa", "temperature", "fahrenheit", -22.0, 86.0},
			CreateContextWithSpan(),
			"2222-bbbb",
			false,
			&Sensor{"2222-bbbb", "1111-aaaa", "temperature", "fahrenheit", -22.0, 86.0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
				tracer: mocktracer.New(),
			}

			// Mock SQL query
			expect := mock.ExpectQuery(`SELECT id, site_id, name, unit, min_safe, max_safe FROM sensors`)
			if tc.dbError != nil {
				expect.WillReturnError(tc.dbError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "site_id", "name", "unit", "min_safe", "max_safe"})
				rows.AddRow(tc.dbRow...)
				expect.WillReturnRows(rows)
			}

			sensor, err := m.Get(tc.context, tc.id)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, sensor)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSensor.ID, sensor.ID)
				assert.Equal(t, tc.expectedSensor.SiteID, sensor.SiteID)
				assert.Equal(t, tc.expectedSensor.Name, sensor.Name)
				assert.Equal(t, tc.expectedSensor.Unit, sensor.Unit)
				assert.Equal(t, tc.expectedSensor.MinSafe, sensor.MinSafe)
				assert.Equal(t, tc.expectedSensor.MaxSafe, sensor.MaxSafe)
			}
		})
	}
}

func TestSensorManagerUpdate(t *testing.T) {
	tests := []struct {
		name               string
		dbError            error
		dbResult           driver.Result
		context            context.Context
		sensor             Sensor
		expectError        bool
		expectedAffectedNo int
	}{
		{
			"DatabaseError",
			errors.New("db error"),
			nil,
			CreateContextWithSpan(),
			Sensor{"", "", "", "", 0.0, 0.0},
			true,
			0,
		},
		{
			"UpdateSensor",
			nil,
			sqlmock.NewResult(0, 1),
			CreateContextWithSpan(),
			Sensor{"2222-bbbb", "1111-aaaa", "temperature", "fahrenheit", -22.0, 86.0},
			false,
			1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
				tracer: mocktracer.New(),
			}

			// Mock SQL query
			expect := mock.ExpectExec(`UPDATE sensors SET`)
			if tc.dbError != nil {
				expect.WillReturnError(tc.dbError)
			} else {
				expect.WillReturnResult(tc.dbResult)
			}

			n, err := m.Update(tc.context, tc.sensor)

			if tc.expectError {
				assert.Error(t, err)
				assert.Zero(t, n)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAffectedNo, n)
			}
		})
	}
}

func TestSensorManagerDelete(t *testing.T) {
	tests := []struct {
		name        string
		dbError     error
		dbResult    driver.Result
		context     context.Context
		sensorID    string
		expectError bool
	}{
		{
			"DatabaseError",
			errors.New("db error"),
			nil,
			CreateContextWithSpan(),
			"",
			true,
		},
		{
			"DeleteSensor",
			nil,
			sqlmock.NewResult(0, 1),
			CreateContextWithSpan(),
			"2222-bbbb",
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
				tracer: mocktracer.New(),
			}

			// Mock SQL query
			expect := mock.ExpectExec(`DELETE FROM sensors`)
			if tc.dbError != nil {
				expect.WillReturnError(tc.dbError)
			} else {
				expect.WillReturnResult(tc.dbResult)
			}

			err = m.Delete(tc.context, tc.sensorID)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
