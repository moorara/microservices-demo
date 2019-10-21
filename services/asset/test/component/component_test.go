package component

import "github.com/moorara/konfig"

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
	_ = konfig.Pick(&Config)
}
