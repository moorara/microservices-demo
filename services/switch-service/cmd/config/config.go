package config

const (
	defaultServiceName      = "switch-service"
	defaultServiceGRPCPort  = ":4030"
	defaultServiceHTTPPort  = ":4031"
	defaultLogLevel         = "info"
	defaultArangoUser       = "root"
	defaultArangoDatabase   = "switches"
	defaultArangoCollection = "switches"
	defaultJaegerAgentAddr  = "localhost:6831"
	defaultJaegerLogSpans   = false
)

var (
	defaultArangoEndpoints = []string{"http://localhost:8529"}
)

// Config defines the schema for configurations
type Config struct {
	ServiceName      string
	ServiceGRPCPort  string
	ServiceHTTPPort  string
	LogLevel         string
	ArangoEndpoints  []string
	ArangoUser       string
	ArangoPassword   string
	ArangoDatabase   string
	ArangoCollection string
	JaegerAgentAddr  string
	JaegerLogSpans   bool
	CAChainFile      string
	ServerCertFile   string
	ServerKeyFile    string
}

// New creates a new configuration object
func New() Config {
	return Config{
		ServiceName:      defaultServiceName,
		ServiceGRPCPort:  defaultServiceGRPCPort,
		ServiceHTTPPort:  defaultServiceHTTPPort,
		LogLevel:         defaultLogLevel,
		ArangoEndpoints:  defaultArangoEndpoints,
		ArangoUser:       defaultArangoUser,
		ArangoDatabase:   defaultArangoDatabase,
		ArangoCollection: defaultArangoCollection,
		JaegerAgentAddr:  defaultJaegerAgentAddr,
		JaegerLogSpans:   defaultJaegerLogSpans,
	}
}
