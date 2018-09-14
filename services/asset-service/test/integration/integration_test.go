package integration

import (
	"github.com/moorara/goto/config"
)

var Config = struct {
	IntegrationTest   bool
	LogLevel          string
	NatsServers       []string
	NatsUser          string
	NatsPassword      string
	CockroachAddr     string
	CockroachUser     string
	CockroachPassword string
	CockroachDatabase string
}{
	LogLevel:          "info",
	NatsServers:       []string{"nats://localhost:4222"},
	NatsUser:          "client",
	NatsPassword:      "pass",
	CockroachAddr:     "localhost:26257",
	CockroachUser:     "root",
	CockroachPassword: "",
	CockroachDatabase: "assets",
}

func init() {
	config.Pick(&Config)
}
