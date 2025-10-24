package main

import (
	"log"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	container := app.NewContainer(cfg)
	server, err := app.NewServer(container)
	if err != nil {
		log.Fatal(err)
	}

	defer container.Cleanup()

	server.Start()
}
