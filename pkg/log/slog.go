package log

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// Logger wraps slog.Logger to provide skip-aware logging
type Logger struct {
	logger *slog.Logger
}

// Logger methods

func (l *Logger) Info(msg string, args ...any) {
	l.logWithSkip(slog.LevelInfo, 1, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logWithSkip(slog.LevelError, 1, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logWithSkip(slog.LevelWarn, 1, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logWithSkip(slog.LevelDebug, 1, msg, args...)
}

// WithGroup returns a new Logger that starts a group
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{logger: l.logger.WithGroup(name)}
}

// With returns a new Logger with additional attributes
func (l *Logger) With(args ...any) *Logger {
	return &Logger{logger: l.logger.With(args...)}
}

// WithGroup returns a Logger that starts a group
func WithGroup(name string) *Logger {
	return &Logger{logger: defaultLogger.logger.WithGroup(name)}
}

// With returns a Logger with additional attributes
func With(args ...any) *Logger {
	return &Logger{logger: defaultLogger.logger.With(args...)}
}

// Context-aware logging functions

// InfoCtx logs at Info level using logger from context (falls back to default)
func InfoCtx(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Info(msg, args...)
}

// ErrorCtx logs at Error level using logger from context
func ErrorCtx(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Error(msg, args...)
}

// WarnCtx logs at Warn level using logger from context
func WarnCtx(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Warn(msg, args...)
}

// DebugContext logs at Debug level using logger from context
func DebugCtx(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Debug(msg, args...)
}

type contextKey string

const loggerKey contextKey = "logger"

// WithContext stores a logger with additional attributes in context
func WithContext(ctx context.Context, args ...any) context.Context {
	logger := FromContext(ctx).With(args...)
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger from context, or returns default logger
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return defaultLogger
	}
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	return defaultLogger
}

// ToContext stores a logger in the context
func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func (l *Logger) logWithSkip(level slog.Level, skip int, msg string, args ...any) {
	ctx := context.Background()
	if !l.logger.Enabled(ctx, level) {
		return
	}

	// Get caller information
	var pcs [1]uintptr
	runtime.Callers(skip+2, pcs[:])
	fs := runtime.CallersFrames(pcs[:])
	_, _ = fs.Next()

	// Create a new record with proper source
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	// Call handler with the record - this preserves With/WithGroup
	_ = l.logger.Handler().Handle(ctx, r)
}
