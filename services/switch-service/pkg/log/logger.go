package log

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Logger wraps a go-kit Logger
type Logger struct {
	serviceName string
	loggerName  string
	Logger      log.Logger
}

// NewVoidLogger creates a new logger for testing purposes
func NewVoidLogger() *Logger {
	return &Logger{
		Logger: log.NewNopLogger(),
	}
}

// NewLogger creates a new logger
func NewLogger(serviceName, loggerName, logLevel string) *Logger {
	l := log.NewJSONLogger(os.Stdout)
	l = log.With(
		l,
		"service", serviceName,
		"logger", loggerName,
		"timestamp", log.DefaultTimestampUTC,
	)

	switch logLevel {
	case "debug", "DEBUG":
		l = level.NewFilter(l, level.AllowDebug())
	case "info", "INFO":
		l = level.NewFilter(l, level.AllowInfo())
	case "warn", "WARN":
		l = level.NewFilter(l, level.AllowWarn())
	case "error", "ERROR":
		l = level.NewFilter(l, level.AllowError())
	}

	return &Logger{
		serviceName: serviceName,
		loggerName:  loggerName,
		Logger:      l,
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
