package config

const (
	defaultServiceName     = "switch-service"
	defaultServicePort     = ":4030"
	defaultLogLevel        = "info"
	defaultJaegerAgentAddr = "localhost:6831"
	defaultJaegerLogSpans  = false
)

// Config defines the schema for configurations
type Config struct {
	ServiceName     string
	ServicePort     string
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
		ServicePort:     defaultServicePort,
		LogLevel:        defaultLogLevel,
		JaegerAgentAddr: defaultJaegerAgentAddr,
		JaegerLogSpans:  defaultJaegerLogSpans,
	}
}
