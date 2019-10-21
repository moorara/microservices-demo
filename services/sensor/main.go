package main

import (
	"github.com/moorara/konfig"
	"github.com/moorara/microservices-demo/services/sensor/config"
	"github.com/moorara/microservices-demo/services/sensor/server"
)

func main() {
	config := config.New()
	err := konfig.Pick(&config)
	if err != nil {
		panic(err)
	}

	server := server.New(config)
	defer server.Close()

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
