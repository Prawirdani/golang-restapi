package log

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/rs/zerolog"
)

// ===================== Zerolog Adapter =====================

// ZerologAdapter adapts zerolog to the Logger interface
type ZerologAdapter struct {
	l zerolog.Logger
}

func NewZerologAdapter(cfg *config.Config) *ZerologAdapter {
	// Inlining common field keys and format with slog version
	zerolog.TimestampFieldName = "timestamp"
	zerolog.CallerFieldName = "source"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var w io.Writer = os.Stdout
	level := zerolog.InfoLevel
	if !cfg.IsProduction() {
		w = zerolog.ConsoleWriter{NoColor: false, Out: os.Stdout, TimeFormat: time.TimeOnly}
		level = zerolog.DebugLevel
	}

	logger := zerolog.New(w).
		With().
		Dict("app", zerolog.Dict().Str("name", cfg.App.Name).Str("version", cfg.App.Version).Str("env", string(cfg.App.Environment))).
		Timestamp().
		Caller().
		Logger().
		Level(level)
	return &ZerologAdapter{
		l: logger,
	}
}

// NewJSONZerolog creates a zerolog logger with JSON output for production
func NewJSONZerolog(w io.Writer) *ZerologAdapter {
	if w == nil {
		w = os.Stdout
	}
	logger := zerolog.New(w).With().Timestamp().Caller().Logger()
	return &ZerologAdapter{l: logger}
}

// Debug logs at Debug level
func (z *ZerologAdapter) Debug(msg string, args ...any) {
	event := z.l.Debug()
	addFields(event, args...).Msg(msg)
}

// Info logs at Info level
func (z *ZerologAdapter) Info(msg string, args ...any) {
	event := z.l.Info()
	addFields(event, args...).Msg(msg)
}

// Warn logs at Warn level
func (z *ZerologAdapter) Warn(msg string, args ...any) {
	event := z.l.Warn()
	addFields(event, args...).Msg(msg)
}

// Error logs at Error level
func (z *ZerologAdapter) Error(msg string, args ...any) {
	event := z.l.Error()
	addFields(event, args...).Msg(msg)
}

// DebugCtx logs at Debug level with context
func (z *ZerologAdapter) DebugCtx(ctx context.Context, msg string, args ...any) {
	logger := z.l
	if l := GetFromContext(ctx); l != nil {
		if za, ok := l.(*ZerologAdapter); ok {
			logger = za.l
		}
	}
	event := logger.Debug()
	addFields(event, args...).Msg(msg)
}

// InfoCtx logs at Info level with context
func (z *ZerologAdapter) InfoCtx(ctx context.Context, msg string, args ...any) {
	logger := z.l
	if l := GetFromContext(ctx); l != nil {
		if za, ok := l.(*ZerologAdapter); ok {
			logger = za.l
		}
	}
	event := logger.Info()
	addFields(event, args...).Msg(msg)
}

// WarnCtx logs at Warn level with context
func (z *ZerologAdapter) WarnCtx(ctx context.Context, msg string, args ...any) {
	logger := z.l
	if l := GetFromContext(ctx); l != nil {
		if za, ok := l.(*ZerologAdapter); ok {
			logger = za.l
		}
	}
	event := logger.Warn()
	addFields(event, args...).Msg(msg)
}

// ErrorCtx logs at Error level with context
func (z *ZerologAdapter) ErrorCtx(ctx context.Context, msg string, args ...any) {
	logger := z.l
	if l := GetFromContext(ctx); l != nil {
		if za, ok := l.(*ZerologAdapter); ok {
			logger = za.l
		}
	}
	event := logger.Error()
	addFields(event, args...).Msg(msg)
}

// With returns a new logger with additional fields
func (z *ZerologAdapter) With(args ...any) Logger {
	ctx := z.l.With()
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			if key, ok := args[i].(string); ok {
				ctx = addContextField(ctx, key, args[i+1])
			}
		}
	}
	return &ZerologAdapter{l: ctx.Logger()}
}

// Helper functions

func addFields(event *zerolog.Event, args ...any) *zerolog.Event {
	// Skip 2 Frame:
	// 1. This method
	// 2. Wrapper (Info, InfoCtx, etc...)
	event.CallerSkipFrame(2)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			if key, ok := args[i].(string); ok {
				event = addEventField(event, key, args[i+1])
			}
		}
	}
	return event
}

func addContextField(ctx zerolog.Context, key string, value any) zerolog.Context {
	switch v := value.(type) {
	case string:
		return ctx.Str(key, v)
	case int:
		return ctx.Int(key, v)
	case int64:
		return ctx.Int64(key, v)
	case int32:
		return ctx.Int32(key, v)
	case float64:
		return ctx.Float64(key, v)
	case float32:
		return ctx.Float32(key, v)
	case bool:
		return ctx.Bool(key, v)
	case error:
		return ctx.AnErr(key, v)
	case []byte:
		return ctx.Bytes(key, v)
	default:
		return ctx.Interface(key, v)
	}
}

func addEventField(event *zerolog.Event, key string, value any) *zerolog.Event {
	switch v := value.(type) {
	case string:
		return event.Str(key, v)
	case int:
		return event.Int(key, v)
	case int64:
		return event.Int64(key, v)
	case int32:
		return event.Int32(key, v)
	case float64:
		return event.Float64(key, v)
	case float32:
		return event.Float32(key, v)
	case bool:
		return event.Bool(key, v)
	case error:
		return event.Err(v)
	case []byte:
		return event.Bytes(key, v)
	default:
		return event.Interface(key, v)
	}
}
