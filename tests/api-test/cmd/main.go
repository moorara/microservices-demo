package main

import (
	"github.com/moorara/microservices-demo/tests/api-test/cmd/config"
	"github.com/moorara/microservices-demo/tests/api-test/cmd/version"
	"github.com/moorara/microservices-demo/tests/api-test/pkg/log"
)

func main() {
	logger := log.NewLogger(config.Config.Name, config.Config.LogLevel)

	logger.Info(
		"version", version.Version,
		"revision", version.Revision,
		"branch", version.Branch,
		"buildTime", version.BuildTime,
		"message", "API tests started.",
	)
}
