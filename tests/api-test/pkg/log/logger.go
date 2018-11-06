package log

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Logger wraps a go-kit Logger
type Logger struct {
	Name   string
	Logger log.Logger
}

// NewNopLogger creates a new logger for testing purposes
func NewNopLogger() *Logger {
	return &Logger{
		Logger: log.NewNopLogger(),
	}
}

// NewLogger creates a new logger
func NewLogger(name, logLevel string) *Logger {
	l := log.NewJSONLogger(os.Stdout)
	l = log.With(
		l,
		"logger", name,
		"timestamp", log.DefaultTimestampUTC,
	)

	switch logLevel {
	case "DEBUG", "Debug", "debug":
		l = level.NewFilter(l, level.AllowDebug())
	case "INFO", "Info", "info":
		l = level.NewFilter(l, level.AllowInfo())
	case "WARN", "Warn", "warn":
		l = level.NewFilter(l, level.AllowWarn())
	case "ERROR", "Error", "error":
		l = level.NewFilter(l, level.AllowError())
	}

	return &Logger{
		Name:   name,
		Logger: l,
	}
}

// Debug logs in debug level
func (l *Logger) Debug(kv ...interface{}) error {
	return level.Debug(l.Logger).Log(kv...)
}

// Info logs in info level
func (l *Logger) Info(kv ...interface{}) error {
	return level.Info(l.Logger).Log(kv...)
}

// Warn logs in warn level
func (l *Logger) Warn(kv ...interface{}) error {
	return level.Warn(l.Logger).Log(kv...)
}

// Error logs in error level
func (l *Logger) Error(kv ...interface{}) error {
	return level.Error(l.Logger).Log(kv...)
}
