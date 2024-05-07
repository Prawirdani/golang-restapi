package main

import (
	"log"
	"log/slog"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal(err)
	}

	app.InitLogger(cfg.App)

	dbPool, err := database.NewPGConnection(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()
	slog.Info("PostgreSQL DB Connection Established")

	router, err := app.InitMainRouter(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Bootstrap(app.Configuration{
		MainRouter: router,
		DBPool:     dbPool,
		Config:     cfg,
	})

	server := app.NewServer(cfg.App, router)
	server.Start()
}
