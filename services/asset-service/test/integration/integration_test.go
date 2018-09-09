package integration

import (
	"github.com/moorara/goto/config"
)

var Config = struct {
	IntegrationTest bool
	NatsServers     []string
	NatsUser        string
	NatsPassword    string
}{
	NatsServers:  []string{"nats://localhost:4222"},
	NatsUser:     "nats_client",
	NatsPassword: "password?!",
}

func init() {
	config.Pick(&Config)
}
