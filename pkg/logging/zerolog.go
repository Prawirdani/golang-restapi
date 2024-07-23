package logging

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/rs/zerolog"
)

type zeroLogger struct {
	logger *zerolog.Logger
}

func newZeroLogger(cfg *config.Config) Logger {
	var w io.Writer = os.Stdout
	var level zerolog.Level = zerolog.DebugLevel

	initTime := time.Now()

	// Write logs into file if in production mode
	if cfg.IsProduction() {
		level = zerolog.InfoLevel

		// Create log directory if not exist
		err := os.MkdirAll(cfg.App.LogPath, 0755)
		if err != nil {
			panic(err)
		}

		// Create log file
		filename := fmt.Sprintf("%s%v.%s", cfg.App.LogPath, initTime.Format("2006-01-02 15:04:05"), "log")

		logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("Failed to open log file")
		}
		w = logFile
	}

	l := zerolog.New(w).With().
		Dict("app", zerolog.Dict().
			Str("name", cfg.App.Name).
			Str("version", cfg.App.Version),
		).
		Timestamp().
		Logger().
		Level(level)

	return &zeroLogger{
		logger: &l,
	}
}

func (zl *zeroLogger) Info(cat Category, caller string, message string) {
	zl.logger.Info().
		Str("category", cat.String()).
		Str("caller", caller).
		Msg(message)
}

func (zl *zeroLogger) Debug(cat Category, caller string, message string) {
	zl.logger.Debug().
		Str("category", cat.String()).
		Str("caller", caller).
		Msg(message)
}

func (zl *zeroLogger) Error(cat Category, caller string, message string) {
	zl.logger.Error().
		Str("category", cat.String()).
		Str("caller", caller).
		Msg(message)
}

func (zl *zeroLogger) Fatal(cat Category, caller string, message string) {
	zl.logger.Fatal().
		Str("category", cat.String()).
		Str("caller", caller).
		Msg(message)
}
