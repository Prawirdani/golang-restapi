package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	cfg := app.NewConfig()
	app.InitLogger(cfg)
	jsonHandler := app.NewJsonHandler()

	dbpool := app.NewPGPool(cfg)
	router := app.InitMainRouter(cfg, jsonHandler)

	validator := validator.New()
	bootstrap := app.Configuration{
		MainRouter: router,
		Config:     cfg,
		DBPool:     dbpool,
		JSON:       jsonHandler,
		Validator:  validator,
	}
	app.Bootstrap(&bootstrap)

	server := app.NewServer(cfg, router)
	server.Start()
}
