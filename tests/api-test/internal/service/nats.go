package service

type (
	// NATSPlan is for testing a service with NATS messaging API
	NATSPlan struct {
		Name         string     `yaml:"name"`
		NATSServers  []string   `yaml:"nats_servers,flow"`
		NATSUser     string     `yaml:"nats_user"`
		NATSPassword string     `yaml:"nats_password"`
		Tests        []NATSTest `yaml:"tests"`
	}

	// NATSTest defines test spec for a NATS service
	NATSTest struct {
	}
)
