package util

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLogger creates a new logger based on go-kit logger
func NewLogger(logLevel, serviceName, loggerName string) log.Logger {
	logger := log.NewJSONLogger(os.Stdout)
	logger = level.NewInjector(logger, level.InfoValue()) // default level
	logger = log.With(
		logger,
		"timestamp", log.DefaultTimestampUTC,
		"service", serviceName,
		"logger", loggerName,
	)

	switch logLevel {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	}

	return logger
}
