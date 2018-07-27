package main

import (
	gotoConfig "github.com/moorara/goto/config"
	"github.com/moorara/microservices-demo/services/switch-service/cmd/config"
	"github.com/moorara/microservices-demo/services/switch-service/internal/server"
)

func main() {
	config := config.New()
	gotoConfig.Pick(&config)

	server, err := server.New(config)
	if err != nil {
		panic(err)
	}

	server.Start()
}
