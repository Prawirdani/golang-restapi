package log

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/prawirdani/golang-restapi/config"
)

var defaultLogger *Logger

func init() {
	defaultLogger = &Logger{logger: slog.Default()}
}

// Info logs at Info level
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Error logs at Error level
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// Warn logs at Warn level
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Debug logs at Debug level
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// InitLogger initializes the global logger
func InitLogger(cfg config.Config) {
	var handler slog.Handler

	if !cfg.IsProduction() {
		handler = NewPrettyHandler(&slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		opts := &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					a.Key = "message"
				}
				if a.Key == slog.TimeKey {
					a.Key = "timestamp"
				}
				if a.Key == slog.SourceKey {
					if source, ok := a.Value.Any().(*slog.Source); ok {
						a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
					}
				}
				return a
			},
		}
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler).With(
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
	)

	defaultLogger = &Logger{logger: logger}
}
