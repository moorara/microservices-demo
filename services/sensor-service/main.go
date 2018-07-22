package main

import (
	gotoConfig "github.com/moorara/goto/config"
	"github.com/moorara/microservices-demo/services/sensor-service/config"
	"github.com/moorara/microservices-demo/services/sensor-service/server"
)

func main() {
	config := config.New()
	gotoConfig.Pick(&config)
	server := server.New(config)
	defer server.Close()

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
