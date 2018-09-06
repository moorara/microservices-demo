package integration

import (
	"github.com/moorara/goto/config"
)

var Config = struct {
	IntegrationTest bool
	ArangoHTTPAddr  string
	ArangoEndpoints []string
	ArangoUser      string
	ArangoPassword  string
	NatsServers     []string
	NatsUser        string
	NatsPassword    string
}{
	ArangoHTTPAddr:  "http://localhost:8529",
	ArangoEndpoints: []string{"tcp://localhost:8529"},
	ArangoUser:      "root",
	ArangoPassword:  "password?!",
	NatsServers:     []string{"nats://localhost:4222"},
	NatsUser:        "nats_client",
	NatsPassword:    "password?!",
}

func init() {
	config.Pick(&Config)
}
