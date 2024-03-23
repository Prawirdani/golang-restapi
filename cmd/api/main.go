package main

import (
	"log"
	"log/slog"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	viper := config.LoadConfig("./config")
	cfg := config.ParseConfig(viper)
	app.InitLogger(cfg.App)

	dbPool, err := database.NewPGConnection(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("PostgreSQL DB Connection Established")

	router := app.InitMainRouter(*cfg)

	app.Bootstrap(&app.Configuration{
		MainRouter: router,
		DBPool:     dbPool,
		Config:     cfg,
	})

	server := app.NewServer(cfg.App, router)
	server.Start()
}
