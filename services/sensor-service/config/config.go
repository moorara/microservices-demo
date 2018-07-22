package config

const (
	dbName                 = "sensors"
	dbOpts                 = "?sslmode=disable"
	defaultServiceName     = "sensor-service"
	defaultServicePort     = ":4020"
	defaultPostgresURL     = "postgres://root@localhost"
	defaultLogLevel        = "info"
	defaultJaegerAgentAddr = "localhost:6831"
	defaultJaegerLogSpans  = false
)

// Config defines the schema for configurations
type Config struct {
	ServiceName     string
	ServicePort     string
	PostgresURL     string
	LogLevel        string
	JaegerAgentAddr string
	JaegerLogSpans  bool
}

// GetFullPostgresURL return the full Postgres URL including database name and options
func (c *Config) GetFullPostgresURL() string {
	return c.PostgresURL + "/" + dbName + dbOpts
}

// New creates a new configuration object
func New() Config {
	return Config{
		ServiceName:     defaultServiceName,
		ServicePort:     defaultServicePort,
		PostgresURL:     defaultPostgresURL,
		LogLevel:        defaultLogLevel,
		JaegerAgentAddr: defaultJaegerAgentAddr,
		JaegerLogSpans:  defaultJaegerLogSpans,
	}
}
