package config

const (
	defaultLogLevel         = "info"
	defaultServiceName      = "sensor-service"
	defaultServicePort      = ":4020"
	defaultPostgresHost     = "localhost"
	defaultPostgresPort     = "5432"
	defaultPostgresDatabase = "sensors"
	defaultPostgresUsername = "root"
	defaultPostgresPassword = ""
	defaultJaegerAgentAddr  = "localhost:6831"
	defaultJaegerLogSpans   = false
)

// Config defines the schema for configurations
type Config struct {
	LogLevel         string
	ServiceName      string
	ServicePort      string
	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string
	PostgresUsername string
	PostgresPassword string
	JaegerAgentAddr  string
	JaegerLogSpans   bool
}

// New creates a new configuration object
func New() Config {
	return Config{
		LogLevel:         defaultLogLevel,
		ServiceName:      defaultServiceName,
		ServicePort:      defaultServicePort,
		PostgresHost:     defaultPostgresHost,
		PostgresPort:     defaultPostgresPort,
		PostgresDatabase: defaultPostgresDatabase,
		PostgresUsername: defaultPostgresUsername,
		PostgresPassword: defaultPostgresPassword,
		JaegerAgentAddr:  defaultJaegerAgentAddr,
		JaegerLogSpans:   defaultJaegerLogSpans,
	}
}
