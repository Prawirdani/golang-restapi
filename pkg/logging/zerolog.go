package logging

import (
	"os"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/rs/zerolog"
)

type zeroLogger struct {
	logger *zerolog.Logger
	// logFile *os.File
}

func newZeroLogger(cfg *config.Config) *zeroLogger {
	level := zerolog.DebugLevel
	if cfg.IsProduction() {
		level = zerolog.InfoLevel
	}

	l := zerolog.New(os.Stdout).With().
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
