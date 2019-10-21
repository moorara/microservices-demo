package integration

import "github.com/moorara/konfig"

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
	_ = konfig.Pick(&Config)
}
