package config

import (
	"github.com/moorara/goto/config"
)

const (
	defaultLogLevel        = "info"
	defaultServiceName     = "asset-service"
	defaultServicePort     = ":4040"
	defaultArangoUser      = "root"
	defaultArangoPassword  = ""
	defaultNatsUser        = ""
	defaultNatsPassword    = ""
	defaultJaegerAgentAddr = "localhost:6831"
	defaultJaegerLogSpans  = false
)

var (
	defaultArangoEndpoints = []string{"tcp://localhost:8529"}
	defaultNatsServers     = []string{"nats://localhost:4222"}
)

// Config defines the configuration values
var Config = struct {
	LogLevel        string
	ServiceName     string
	ServicePort     string
	ArangoEndpoints []string
	ArangoUser      string
	ArangoPassword  string
	NatsServers     []string
	NatsUser        string
	NatsPassword    string
	JaegerAgentAddr string
	JaegerLogSpans  bool
}{
	LogLevel:        defaultLogLevel,
	ServiceName:     defaultServiceName,
	ServicePort:     defaultServicePort,
	ArangoEndpoints: defaultArangoEndpoints,
	ArangoUser:      defaultArangoUser,
	ArangoPassword:  defaultArangoPassword,
	NatsServers:     defaultNatsServers,
	NatsUser:        defaultNatsUser,
	NatsPassword:    defaultNatsPassword,
	JaegerAgentAddr: defaultJaegerAgentAddr,
	JaegerLogSpans:  defaultJaegerLogSpans,
}

func init() {
	config.Pick(&Config)
}
