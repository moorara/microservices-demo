package component

import (
	"github.com/moorara/goto/config"
)

var Config = struct {
	ComponentTest bool
	ServiceURL    string
	NatsServers   []string
	NatsUser      string
	NatsPassword  string
}{
	ServiceURL:   "http://localhost:4040",
	NatsServers:  []string{"nats://localhost:4222"},
	NatsUser:     "client",
	NatsPassword: "pass",
}

func init() {
	config.Pick(&Config)
}
