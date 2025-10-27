package main

import (
	"os"

	stdlog "log"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/app"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		stdlog.Fatal("failed to load config", err)
	}
	log.InitLogger(*cfg)

	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Error("failed to create container", "error", err.Error())
		os.Exit(1)
	}
	server, err := app.NewServer(container)
	if err != nil {
		log.Error("failed to create server", "error", err.Error())
		os.Exit(1)
	}

	defer container.Cleanup()

	server.Start()
}
