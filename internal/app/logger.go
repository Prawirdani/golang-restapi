package app

import (
	"log"
	"os"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
)

func InitLogger(config *viper.Viper) {
	handler := new(slog.HandlerOptions)

	currentEnv := strings.ToUpper(config.GetString("app.env"))
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
