package config

const (
	dbName             = "sensors"
	dbOpts             = "?sslmode=disable"
	defaultServiceName = "sensor-service"
	defaultServicePort = ":4020"
	defaultPostgresURL = "postgres://root@localhost"
	defaultLogLevel    = "info"
)

// Spec represents configuration specifications
type Spec struct {
	ServiceName string
	ServicePort string
	PostgresURL string
	LogLevel    string
}

// GetFullPostgresURL return the full Postgres URL including database name and options
func (s *Spec) GetFullPostgresURL() string {
	return s.PostgresURL + "/" + dbName + dbOpts
}

// Config is the configuration object
var Config = Spec{
	ServiceName: defaultServiceName,
	ServicePort: defaultServicePort,
	PostgresURL: defaultPostgresURL,
	LogLevel:    defaultLogLevel,
}
