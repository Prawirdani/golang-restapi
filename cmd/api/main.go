package main

import (
	"log"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/app"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(cfg)
	dbPool, err := database.NewPGConnection(cfg)
	if err != nil {
		logger.Fatal(logging.Postgres, "main.NewPGConnection", err.Error())
	}
	defer dbPool.Close()
	logger.Info(logging.Startup, "main", "Postgres connection established")

	server, err := app.InitServer(cfg, logger, dbPool)
	if err != nil {
		logger.Fatal(logging.Startup, "main.InitServer", err.Error())
	}

	server.Start()
}
