package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	CloseCalled bool
	CloseError  error

	ExecContextCalled bool
	ExecContextResult sql.Result
	ExecContextError  error

	QueryContextCalled bool
	QueryContextRows   *sql.Rows
	QueryContextError  error

	QueryRowContextCalled bool
	QueryRowContextRow    *sql.Row
}

func (mdb *mockDB) Close() error {
	mdb.CloseCalled = true
	return mdb.CloseError
}

func (mdb *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mdb.ExecContextCalled = true
	return mdb.ExecContextResult, mdb.ExecContextError
}

func (mdb *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	mdb.QueryContextCalled = true
	return mdb.QueryContextRows, mdb.QueryContextError
}

func (mdb *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	mdb.QueryRowContextCalled = true
	return mdb.QueryRowContextRow
}

type mockResult struct {
	LastInsertIDCalled bool
	LastInsertIDResult int64
	LastInsertIDError  error

	RowsAffectedCalled bool
	RowsAffectedResult int64
	RowsAffectedError  error
}

func (mr *mockResult) LastInsertId() (int64, error) {
	mr.LastInsertIDCalled = true
	return mr.LastInsertIDResult, mr.LastInsertIDError
}

func (mr *mockResult) RowsAffected() (int64, error) {
	mr.RowsAffectedCalled = true
	return mr.RowsAffectedResult, mr.RowsAffectedError
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
			m := NewSensorManager(db, logger)

			assert.NotNil(t, m)
		})
	}
}

func TestSensorManagerCreate(t *testing.T) {
	tests := []struct {
		name             string
		execContextError error
		context          context.Context
		sensorSiteID     string
		sensorName       string
		sensorUnit       string
		sensorMinSafe    float64
		sensorMaxSafe    float64
		expectError      bool
	}{
		{
			"DatabaseError",
			errors.New("error"),
			context.Background(),
			"", "", "", 0.0, 0.0,
			true,
		},
		{
			"CreateSensor",
			nil,
			context.Background(),
			"1111-aaaa", "temperature", "celsius", -30.0, 30.0,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := &mockDB{
				ExecContextResult: nil,
				ExecContextError:  tc.execContextError,
			}

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
			}

			sensor, err := m.Create(tc.context, tc.sensorSiteID, tc.sensorName, tc.sensorUnit, tc.sensorMinSafe, tc.sensorMaxSafe)

			assert.True(t, db.ExecContextCalled)
			if tc.expectError {
				assert.Error(t, err)
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

func TestSensorManagerUpdate(t *testing.T) {
	tests := []struct {
		name               string
		execContextResult  *mockResult
		execContextError   error
		context            context.Context
		sensor             Sensor
		expectError        bool
		expectedAffectedNo int
	}{
		{
			"DatabaseError",
			&mockResult{},
			errors.New("error"),
			context.Background(),
			Sensor{"", "", "", "", 0.0, 0.0},
			true,
			0,
		},
		{
			"ResultError",
			&mockResult{
				RowsAffectedResult: 0,
				RowsAffectedError:  errors.New("error"),
			},
			nil,
			context.Background(),
			Sensor{"", "", "", "", 0.0, 0.0},
			true,
			0,
		},
		{
			"UpdateSensor",
			&mockResult{
				RowsAffectedResult: 1,
				RowsAffectedError:  nil,
			},
			nil,
			context.Background(),
			Sensor{"2222-bbbb", "1111-aaaa", "temperature", "fahrenheit", -22.0, 86.0},
			false,
			1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := &mockDB{
				ExecContextResult: tc.execContextResult,
				ExecContextError:  tc.execContextError,
			}

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
			}

			n, err := m.Update(tc.context, tc.sensor)

			assert.True(t, db.ExecContextCalled)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAffectedNo, n)
			}
		})
	}
}

func TestSensorManagerDelete(t *testing.T) {
	tests := []struct {
		name             string
		execContextError error
		context          context.Context
		sensorID         string
		expectError      bool
	}{
		{
			"DatabaseError",
			errors.New("error"),
			context.Background(),
			"",
			true,
		},
		{
			"DeleteSensor",
			nil,
			context.Background(),
			"2222-bbbb",
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := &mockDB{
				ExecContextResult: nil,
				ExecContextError:  tc.execContextError,
			}

			m := &postgresSensorManager{
				db:     db,
				logger: log.NewNopLogger(),
			}

			err := m.Delete(tc.context, tc.sensorID)

			assert.True(t, db.ExecContextCalled)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
