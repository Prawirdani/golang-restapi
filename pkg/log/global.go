package log

import (
	"context"
	"log/slog"
	"os"
)

// global is the package-level default logger instance.
// It defaults to a slog-based JSON logger that writes to os.Stdout.
//
// This instance is used by all global logging functions (e.g., Info(), Error())
// when no explicit logger or context-bound logger is available./ Default to slog
var global Logger = defaultAdapter()

func defaultAdapter() *SlogAdapter {
	return &SlogAdapter{l: slog.New(slog.NewTextHandler(os.Stdout, nil))}
}

// SetLogger replaces the global logger used by all package-level logging calls.
//
// Passing a nil logger will cause a panic. This function should typically be called
// once during application initialization to configure the desired logging backend
// (for example, a development pretty printer or a production JSON logger).
func SetLogger(l Logger) {
	if l == nil {
		panic("logger: SetLogger called with nil")
	}
	global = l
}

func Debug(msg string, args ...any)            { global.Debug(msg, args...) }
func Info(msg string, args ...any)             { global.Info(msg, args...) }
func Warn(msg string, args ...any)             { global.Warn(msg, args...) }
func Error(msg string, err error, args ...any) { global.Error(msg, err, args...) }

func DebugCtx(ctx context.Context, msg string, args ...any) { global.DebugCtx(ctx, msg, args...) }
func InfoCtx(ctx context.Context, msg string, args ...any)  { global.InfoCtx(ctx, msg, args...) }
func WarnCtx(ctx context.Context, msg string, args ...any)  { global.WarnCtx(ctx, msg, args...) }
func ErrorCtx(ctx context.Context, msg string, err error, args ...any) {
	global.ErrorCtx(ctx, msg, err, args...)
}
