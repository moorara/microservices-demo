package main

import (
	"github.com/moorara/toys/microservices/go-service/config"
	"github.com/moorara/toys/microservices/go-service/server"
)

func main() {
	config := config.GetConfig()
	server := server.New(config)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
