package config

import (
	"github.com/moorara/goto/config"
)

const (
	defaultLogLevel          = "info"
	defaultServiceName       = "asset-service"
	defaultServicePort       = ":4040"
	defaultNatsUser          = ""
	defaultNatsPassword      = ""
	defaultCockroachAddr     = "localhost:26257"
	defaultCockroachUser     = "root"
	defaultCockroachPassword = ""
	defaultJaegerAgentAddr   = "localhost:6831"
	defaultJaegerLogSpans    = false
)

var (
	defaultNatsServers = []string{"nats://localhost:4222"}
)

// Config defines the configuration values
var Config = struct {
	LogLevel          string
	ServiceName       string
	ServicePort       string
	NatsServers       []string
	NatsUser          string
	NatsPassword      string
	CockroachAddr     string
	CockroachUser     string
	CockroachPassword string
	JaegerAgentAddr   string
	JaegerLogSpans    bool
}{
	LogLevel:          defaultLogLevel,
	ServiceName:       defaultServiceName,
	ServicePort:       defaultServicePort,
	NatsServers:       defaultNatsServers,
	NatsUser:          defaultNatsUser,
	NatsPassword:      defaultNatsPassword,
	CockroachAddr:     defaultCockroachAddr,
	CockroachUser:     defaultCockroachUser,
	CockroachPassword: defaultCockroachPassword,
	JaegerAgentAddr:   defaultJaegerAgentAddr,
	JaegerLogSpans:    defaultJaegerLogSpans,
}

func init() {
	config.Pick(&Config)
}
