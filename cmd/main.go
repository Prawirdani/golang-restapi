package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	cfg := app.NewConfig()
	app.InitLogger(cfg)

	dbPool := app.NewPGPool(cfg)
	router := app.InitMainRouter(cfg)

	validator := validator.New()
	bootstrap := app.Configuration{
		MainRouter: router,
		Config:     cfg,
		DBPool:     dbPool,
		Validator:  validator,
	}
	app.Bootstrap(&bootstrap)

	server := app.NewServer(cfg, router)
	server.Start()
}
