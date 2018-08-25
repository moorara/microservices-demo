package config

import (
	"time"
)

const (
	defaultLogLevel    = "info"
	defaultServiceName = "switch-service"

	defaultServiceGRPCPort = ":4030"
	defaultServiceHTTPPort = ":4031"
	defaultServerTimeout   = int64(30 * time.Second)

	defaultArangoUser       = "root"
	defaultArangoDatabase   = "switches"
	defaultArangoCollection = "switches"

	defaultJaegerAgentAddr = "localhost:6831"
	defaultJaegerLogSpans  = false
)

var (
	defaultArangoEndpoints = []string{"tcp://localhost:8529"}
)

// Config defines the schema for configurations
type Config struct {
	LogLevel    string
	ServiceName string

	ServiceGRPCPort string
	ServiceHTTPPort string
	ServerTimeout   int64

	ArangoEndpoints  []string
	ArangoUser       string
	ArangoPassword   string
	ArangoDatabase   string
	ArangoCollection string

	JaegerAgentAddr string
	JaegerLogSpans  bool

	CAChainFile    string
	ServerCertFile string
	ServerKeyFile  string
}

// New creates a new configuration object
func New() Config {
	return Config{
		LogLevel:    defaultLogLevel,
		ServiceName: defaultServiceName,

		ServiceGRPCPort: defaultServiceGRPCPort,
		ServiceHTTPPort: defaultServiceHTTPPort,
		ServerTimeout:   defaultServerTimeout,

		ArangoEndpoints:  defaultArangoEndpoints,
		ArangoUser:       defaultArangoUser,
		ArangoDatabase:   defaultArangoDatabase,
		ArangoCollection: defaultArangoCollection,

		JaegerAgentAddr: defaultJaegerAgentAddr,
		JaegerLogSpans:  defaultJaegerLogSpans,
	}
}
