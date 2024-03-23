package app

import (
	"log"
	"os"
	"strings"

	"log/slog"

	"github.com/prawirdani/golang-restapi/config"
)

func InitLogger(cfg config.AppConfig) {
	handler := new(slog.HandlerOptions)

	currentEnv := strings.ToUpper(cfg.Environment)
	log.Printf("App Version: %s\n", cfg.Version)
	log.Printf("Environtment: %s\n", currentEnv)

	if currentEnv == "PROD" {
		log.Println("Log Level: Info")
		handler.Level = slog.LevelInfo
	} else {
		log.Println("Log Level: Debug")
		handler.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, handler))
	slog.SetDefault(logger)
}
