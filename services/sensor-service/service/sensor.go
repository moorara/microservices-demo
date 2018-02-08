package service

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
)

type (
	// Sensor represents a Sensor in a site
	Sensor struct {
		ID      string  `json:"id"`
		SiteID  string  `json:"siteId"`
		Name    string  `json:"name"`
		Unit    string  `json:"unit"`
		MinSafe float64 `json:"minSafe"`
		MaxSafe float64 `json:"maxSafe"`
	}

	// SensorManager abstracts CRUD operations for Sensor
	SensorManager interface {
		Create(ctx context.Context, siteID string, name, unit string, minSafe, maxSafe float64) (*Sensor, error)
		All(ctx context.Context, siteID string) ([]Sensor, error)
		Get(ctx context.Context, id string) (*Sensor, error)
		Delete(ctx context.Context, id string) error
	}

	postgresSensorManager struct {
		db     DB
		logger log.Logger
	}
)

// NewSensorManager creates a new sensor manager
func NewSensorManager(db DB, logger log.Logger) SensorManager {
	return &postgresSensorManager{
		db:     db,
		logger: logger,
	}
}

func (m *postgresSensorManager) Create(ctx context.Context, siteID string, name, unit string, minSafe, maxSafe float64) (*Sensor, error) {
	sensor := &Sensor{
		ID:      uuid.New().String(),
		SiteID:  siteID,
		Name:    name,
		Unit:    unit,
		MinSafe: minSafe,
		MaxSafe: maxSafe,
	}

	query := `insert into sensors (id, site_id, name, unit, min_safe, max_safe) values ($1, $2, $3, $4, $5, $6)`
	_, err := m.db.ExecContext(ctx, query, sensor.ID, sensor.SiteID, sensor.Name, sensor.Unit, sensor.MinSafe, sensor.MaxSafe)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return nil, err
	}

	return sensor, nil
}

func (m *postgresSensorManager) All(ctx context.Context, siteID string) ([]Sensor, error) {
	sensors := make([]Sensor, 0)

	query := `select id, site_id, name, unit, min_safe, max_safe from sensors where site_id = $1`
	rows, err := m.db.QueryContext(ctx, query, siteID)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return nil, err
	}

	for rows.Next() {
		sensor := Sensor{}
		err := rows.Scan(&sensor.ID, &sensor.SiteID, &sensor.Name, &sensor.Unit, &sensor.MinSafe, &sensor.MaxSafe)
		if err == nil {
			sensors = append(sensors, sensor)
		}
	}

	return sensors, nil
}

func (m *postgresSensorManager) Get(ctx context.Context, id string) (*Sensor, error) {
	sensor := new(Sensor)

	query := `select id, site_id, name, unit, min_safe, max_safe from sensors where id = $1`
	row := m.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&sensor.ID, &sensor.SiteID, &sensor.Name, &sensor.Unit, &sensor.MinSafe, &sensor.MaxSafe)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		level.Error(m.logger).Log("message", err.Error())
		return nil, err
	}

	return sensor, nil
}

func (m *postgresSensorManager) Delete(ctx context.Context, id string) error {
	query := `delete from sensors where id = $1`
	_, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return err
	}

	return nil
}
