package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq" // Required for initialization
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

// NewPostgresDB creates a new DB for PostgreSQL
func NewPostgresDB(logger log.Logger, host, port, database, username, password string) DB {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)
	level.Debug(logger).Log("message", "postgres connection string", "connStr", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return db
}
