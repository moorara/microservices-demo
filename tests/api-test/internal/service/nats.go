package service

type (
	// NATSTest is for testing a service with messaging API
	NATSTest struct {
		Name         string   `yaml:"name"`
		NATSServers  []string `yaml:"nats_servers,flow"`
		NATSUser     string   `yaml:"nats_user"`
		NATSPassword string   `yaml:"nats_password"`
	}
)
