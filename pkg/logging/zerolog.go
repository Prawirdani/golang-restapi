package logging

import (
	"os"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/rs/zerolog"
)

type zeroLogger struct {
	logger *zerolog.Logger
}

func newZeroLogger(cfg *config.Config) *zeroLogger {
	if cfg.IsProduction() {
		l := zerolog.New(os.Stdout).With().
			Dict("APP", zerolog.Dict().
				Str("NAME", cfg.App.Name).
				Str("VERSION", cfg.App.Version),
			).
			Timestamp().
			Logger().
			Level(zerolog.InfoLevel)

		return &zeroLogger{
			logger: &l,
		}

	}

	devLogger := zerolog.New(&zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
		NoColor:    false, // Enable colored output
	}).With().
		Str("APP_NAME", cfg.App.Name).
		Str("APP_VERSION", cfg.App.Version).
		Timestamp().
		Logger().
		Level(zerolog.DebugLevel)

	return &zeroLogger{
		logger: &devLogger,
	}
}

func (zl *zeroLogger) Info(cat Category, caller string, message string) {
	zl.logger.Info().
		Str("CATEGORY", cat.String()).
		Str("CALLER", caller).
		Msg(message)
}

func (zl *zeroLogger) Debug(cat Category, caller string, message string) {
	zl.logger.Debug().
		Str("CATEGORY", cat.String()).
		Str("CALLER", caller).
		Msg(message)
}

func (zl *zeroLogger) Error(cat Category, caller string, message string) {
	zl.logger.Error().
		Str("CATEGORY", cat.String()).
		Str("CALLER", caller).
		Msg(message)
}

func (zl *zeroLogger) Fatal(cat Category, caller string, message string) {
	zl.logger.Fatal().
		Str("CATEGORY", cat.String()).
		Str("CALLER", caller).
		Msg(message)
}
