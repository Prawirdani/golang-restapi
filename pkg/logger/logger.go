package logger

import (
	"log/slog"
	"os"
)

type Layer string

const (
	Handler Layer = "Handler"
	Repo    Layer = "Repository"
	Service Layer = "Service"
	Utility Layer = "Utility"
)

func (s Layer) String() string {
	return string(s)
}

func Init(production bool) {
	handler := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	if production {
		handler.Level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, handler))
	slog.SetDefault(logger)
}

func Info(scope Layer, executor string, message string) {
	slog.Info(message, slog.String("layer", scope.String()), slog.String("executor", executor))
}

func Debug(scope Layer, executor string, message string) {
	slog.Debug(message, slog.String("layer", scope.String()), slog.String("executor", executor))
}

func Error(scope Layer, executor string, message string) {
	slog.Error(message, slog.String("layer", scope.String()), slog.String("executor", executor))
}
