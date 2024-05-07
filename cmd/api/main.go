package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/database"
	"github.com/prawirdani/golang-restapi/internal/app"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Version: %s, Environtment: %s", cfg.App.Version, cfg.App.Environment)

	initAppLogger(cfg.IsProduction())

	dbPool, err := database.NewPGConnection(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()
	slog.Info("PostgreSQL DB Connection Established")

	server, err := app.InitServer(cfg, dbPool)
	if err != nil {
		log.Fatal(err)
	}

	server.Start()
}

func initAppLogger(prodEnv bool) {
	handler := new(slog.HandlerOptions)

	if prodEnv {
		log.Println("Log Level: Info")
		handler.Level = slog.LevelInfo
	} else {
		log.Println("Log Level: Debug")
		handler.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, handler))
	slog.SetDefault(logger)
}
