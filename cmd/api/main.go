package main

import (
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	viper := config.LoadConfig()
	cfg := config.ParseConfig(viper)
	app.InitLogger(cfg.App)

	dbPool := database.NewPGPool(cfg.DB)
	router := app.InitMainRouter(*cfg)

	app.Bootstrap(&app.Configuration{
		MainRouter: router,
		DBPool:     dbPool,
		Config:     cfg,
	})

	server := app.NewServer(cfg.App, router)
	server.Start()
}
