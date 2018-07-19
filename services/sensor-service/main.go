package main

import (
	gotoConfig "github.com/moorara/goto/config"
	"github.com/moorara/microservices-demo/services/sensor-service/config"
	"github.com/moorara/microservices-demo/services/sensor-service/server"
)

func main() {
	gotoConfig.Pick(&config.Config)
	server := server.New(config.Config)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
