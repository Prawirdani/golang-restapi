package main

import (
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	cfg := app.ParseConfig()
	app.InitLogger(cfg)

	dbPool := database.NewPGPool(cfg)
	router := app.InitMainRouter(cfg)

	app.Bootstrap(&app.Configuration{
		MainRouter: router,
		DBPool:     dbPool,
		Config:     cfg,
	})

	server := app.NewServer(cfg, router)
	server.Start()
}
