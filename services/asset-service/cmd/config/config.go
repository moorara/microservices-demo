package config

import (
	"github.com/moorara/goto/config"
)

const (
	defaultServiceName     = "asset-service"
	defaultServicePort     = ":4040"
	defaultLogLevel        = "info"
	defaultJaegerAgentAddr = "localhost:6831"
	defaultJaegerLogSpans  = false
)

// Config defines the configuration values
var Config = struct {
	ServiceName     string
	ServicePort     string
	LogLevel        string
	JaegerAgentAddr string
	JaegerLogSpans  bool
}{
	ServiceName:     defaultServiceName,
	ServicePort:     defaultServicePort,
	LogLevel:        defaultLogLevel,
	JaegerAgentAddr: defaultJaegerAgentAddr,
	JaegerLogSpans:  defaultJaegerLogSpans,
}

func init() {
	config.Pick(&Config)
}
