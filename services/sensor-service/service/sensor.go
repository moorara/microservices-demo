package service

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"

	opentracingLog "github.com/opentracing/opentracing-go/log"
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
		Create(ctx context.Context, siteID, name, unit string, minSafe, maxSafe float64) (*Sensor, error)
		All(ctx context.Context, siteID string) ([]Sensor, error)
		Get(ctx context.Context, id string) (*Sensor, error)
		Update(ctx context.Context, s Sensor) (int, error)
		Delete(ctx context.Context, id string) error
	}

	postgresSensorManager struct {
		db     DB
		logger log.Logger
		tracer opentracing.Tracer
	}
)

// NewSensorManager creates a new sensor manager
func NewSensorManager(db DB, logger log.Logger, tracer opentracing.Tracer) SensorManager {
	return &postgresSensorManager{
		db:     db,
		logger: logger,
		tracer: tracer,
	}
}

func (m *postgresSensorManager) Create(ctx context.Context, siteID, name, unit string, minSafe, maxSafe float64) (*Sensor, error) {
	sensor := &Sensor{
		ID:      uuid.New().String(),
		SiteID:  siteID,
		Name:    name,
		Unit:    unit,
		MinSafe: minSafe,
		MaxSafe: maxSafe,
	}

	parentSpan := opentracing.SpanFromContext(ctx)
	span := m.tracer.StartSpan("insert-record", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	query := `insert into sensors (id, site_id, name, unit, min_safe, max_safe) values ($1, $2, $3, $4, $5, $6)`
	_, err := m.db.ExecContext(ctx, query, sensor.ID, sensor.SiteID, sensor.Name, sensor.Unit, sensor.MinSafe, sensor.MaxSafe)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return nil, err
	}

	span.SetTag("db.type", "sql")
	span.LogFields(
		opentracingLog.String("event", "insert-record"),
		opentracingLog.String("db.statement", query),
	)

	return sensor, nil
}

func (m *postgresSensorManager) All(ctx context.Context, siteID string) ([]Sensor, error) {
	sensors := make([]Sensor, 0)

	parentSpan := opentracing.SpanFromContext(ctx)
	span := m.tracer.StartSpan("select-records", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	query := `select id, site_id, name, unit, min_safe, max_safe from sensors where site_id = $1`
	rows, err := m.db.QueryContext(ctx, query, siteID)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return nil, err
	}

	span.SetTag("db.type", "sql")
	span.LogFields(
		opentracingLog.String("event", "select-records"),
		opentracingLog.String("db.statement", query),
	)

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

	parentSpan := opentracing.SpanFromContext(ctx)
	span := m.tracer.StartSpan("select-record", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

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

	span.SetTag("db.type", "sql")
	span.LogFields(
		opentracingLog.String("event", "select-record"),
		opentracingLog.String("db.statement", query),
	)

	return sensor, nil
}

func (m *postgresSensorManager) Update(ctx context.Context, s Sensor) (int, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := m.tracer.StartSpan("update-record", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	query := `update sensors set site_id = $2, name = $3, unit = $4, min_safe = $5, max_safe = $6 where id = $1`
	res, err := m.db.ExecContext(ctx, query, s.ID, s.SiteID, s.Name, s.Unit, s.MinSafe, s.MaxSafe)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return 0, err
	}

	span.SetTag("db.type", "sql")
	span.LogFields(
		opentracingLog.String("event", "update-record"),
		opentracingLog.String("db.statement", query),
	)

	n, err := res.RowsAffected()
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return 0, err
	}

	return int(n), nil
}

func (m *postgresSensorManager) Delete(ctx context.Context, id string) error {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := m.tracer.StartSpan("delete-record", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	query := `delete from sensors where id = $1`
	_, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		level.Error(m.logger).Log("message", err.Error())
		return err
	}

	span.SetTag("db.type", "sql")
	span.LogFields(
		opentracingLog.String("event", "delete-record"),
		opentracingLog.String("db.statement", query),
	)

	return nil
}
