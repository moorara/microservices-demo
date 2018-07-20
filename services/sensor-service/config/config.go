package config

import "fmt"

const (
	dbName                        = "sensors"
	dbOpts                        = "?sslmode=disable"
	defaultServiceName            = "sensor-service"
	defaultServicePort            = ":4020"
	defaultPostgresURL            = "postgres://root@localhost"
	defaultLogLevel               = "info"
	defaultJaegerAgentHost        = "localhost"
	defaultJaegerAgentPort        = 6831
	defaultJaegerReporterLogSpans = false
)

// Config defines the schema for configurations
type Config struct {
	ServiceName            string
	ServicePort            string
	PostgresURL            string
	LogLevel               string
	JaegerAgentHost        string
	JaegerAgentPort        int
	JaegerReporterLogSpans bool
}

// GetFullPostgresURL return the full Postgres URL including database name and options
func (c *Config) GetFullPostgresURL() string {
	return c.PostgresURL + "/" + dbName + dbOpts
}

// GetJaegerAgentURL return the full Jaeger URL including host and port
func (c *Config) GetJaegerAgentURL() string {
	return fmt.Sprintf("%s:%d", c.JaegerAgentHost, c.JaegerAgentPort)
}

// New creates a new configuration object
func New() Config {
	return Config{
		ServiceName:            defaultServiceName,
		ServicePort:            defaultServicePort,
		PostgresURL:            defaultPostgresURL,
		LogLevel:               defaultLogLevel,
		JaegerAgentHost:        defaultJaegerAgentHost,
		JaegerAgentPort:        defaultJaegerAgentPort,
		JaegerReporterLogSpans: defaultJaegerReporterLogSpans,
	}
}
