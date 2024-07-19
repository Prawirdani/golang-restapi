package main

import (
	"log"
	"log/slog"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/api"
	"github.com/prawirdani/golang-restapi/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal(err)
	}
	logger.Init(cfg.IsProduction())
	log.Printf("Version: %s, Environtment: %s", cfg.App.Version, cfg.App.Environment)

	dbPool, err := database.NewPGConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()
	slog.Info("PostgreSQL DB Connection Established")

	server, err := api.InitServer(cfg, dbPool)
	if err != nil {
		log.Fatal(err)
	}

	server.Start()
}
