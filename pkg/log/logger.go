package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/prawirdani/golang-restapi/config"
)

var logger *slog.Logger = slog.Default()

func Info(msg string, args ...any) {
	WithSkip(slog.LevelInfo, 1, msg, args...)
}

func Error(msg string, args ...any) {
	WithSkip(slog.LevelError, 1, msg, args...)
}

func Warn(msg string, args ...any) {
	WithSkip(slog.LevelWarn, 1, msg, args...)
}

func Debug(msg string, args ...any) {
	WithSkip(slog.LevelDebug, 1, msg, args...)
}

func WithSkip(level slog.Level, skip int, msg string, args ...any) {
	if !logger.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	// Skip: runtime.Callers + logWithSkip + Info/Error/etc + actual caller
	runtime.Callers(skip+2, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	_ = logger.Handler().Handle(context.Background(), r)
}

func InitLogger(cfg config.Config) {
	level := slog.LevelDebug

	if cfg.IsProduction() {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // Adds file:line to logs
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Rename "msg" to "message" for better Loki compatibility
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			// Rename "time" to match standard formats
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
			}
			// Format source as "file:line" instead of object
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
				}
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	l := slog.New(handler).With(
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
	)

	logger = l
}
