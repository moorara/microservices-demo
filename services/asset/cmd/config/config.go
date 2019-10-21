package config

import "github.com/moorara/konfig"

const (
	defaultLogLevel          = "info"
	defaultServiceName       = "asset-service"
	defaultServicePort       = ":4040"
	defaultNatsUser          = "client"
	defaultNatsPassword      = "pass"
	defaultCockroachAddr     = "localhost:26257"
	defaultCockroachUser     = "root"
	defaultCockroachPassword = ""
	defaultCockroachDatabase = "assets"
	defaultJaegerAgentAddr   = "localhost:6831"
	defaultJaegerLogSpans    = false
)

var (
	defaultNatsServers = []string{"nats://localhost:4222"}
)

// Global defines the configuration values
var Global = struct {
	LogLevel          string
	ServiceName       string
	ServicePort       string
	NatsServers       []string
	NatsUser          string
	NatsPassword      string
	CockroachAddr     string
	CockroachUser     string
	CockroachPassword string
	CockroachDatabase string
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
	CockroachDatabase: defaultCockroachDatabase,
	JaegerAgentAddr:   defaultJaegerAgentAddr,
	JaegerLogSpans:    defaultJaegerLogSpans,
}

func init() {
	_ = konfig.Pick(&Global)
}
