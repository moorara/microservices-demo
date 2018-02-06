package service

import (
	"context"
	"database/sql"

	// Required for initialization
	_ "github.com/lib/pq"
)

type (
	// DB is the interface for a sql database
	DB interface {
		Close() error
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}
)

// NewPostgresDB creates a new DB for Postgres
func NewPostgresDB(postgresURL string) DB {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		panic(err)
	}

	return db
}
