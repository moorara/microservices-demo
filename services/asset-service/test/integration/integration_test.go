package integration

import (
	"github.com/moorara/goto/config"
)

var Config = struct {
	IntegrationTest   bool
	NatsServers       []string
	NatsUser          string
	NatsPassword      string
	CockroachAddr     string
	CockroachUser     string
	CockroachPassword string
}{
	NatsServers:       []string{"nats://localhost:4222"},
	NatsUser:          "nats_client",
	NatsPassword:      "password?!",
	CockroachAddr:     "localhost:26257",
	CockroachUser:     "service",
	CockroachPassword: "",
}

func init() {
	config.Pick(&Config)
}
