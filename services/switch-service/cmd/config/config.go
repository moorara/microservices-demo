package config

const (
	defaultServiceName     = "switch-service"
	defaultServiceGRPCPort = ":4030"
	defaultServiceHTTPPort = ":4031"
	defaultLogLevel        = "info"
	defaultJaegerAgentAddr = "localhost:6831"
	defaultJaegerLogSpans  = false
)

// Config defines the schema for configurations
type Config struct {
	ServiceName     string
	ServiceGRPCPort string
	ServiceHTTPPort string
	LogLevel        string
	JaegerAgentAddr string
	JaegerLogSpans  bool
	CAChainFile     string
	ServerCertFile  string
	ServerKeyFile   string
}

// New creates a new configuration object
func New() Config {
	return Config{
		ServiceName:     defaultServiceName,
		ServiceGRPCPort: defaultServiceGRPCPort,
		ServiceHTTPPort: defaultServiceHTTPPort,
		LogLevel:        defaultLogLevel,
		JaegerAgentAddr: defaultJaegerAgentAddr,
		JaegerLogSpans:  defaultJaegerLogSpans,
	}
}
