package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/prawirdani/golang-restapi/config"
)

// ===================== Slog Adapter =====================

type SlogAdapter struct {
	l *slog.Logger
}

func NewSlogAdapter(cfg *config.Config) *SlogAdapter {
	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			if a.Key == slog.LevelKey {
				a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
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
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})

	if cfg.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	}

	l := slog.New(handler).With(
		slog.Group("app",
			slog.String("name", cfg.App.Name),
			slog.String("version", cfg.App.Version),
			slog.String("env", string(cfg.App.Environment)),
		),
	)

	return &SlogAdapter{l: l}
}

// Addtional 2 skips to capture the correct caller frame:
// Frame 3: this function
// Frame 4: wrapper (InfoCtx, DebugCtx ...)

func (s *SlogAdapter) Debug(msg string, args ...any) {
	s.logWithSkip(context.Background(), slog.LevelDebug, 2, msg, args...)
}

func (s *SlogAdapter) Info(msg string, args ...any) {
	s.logWithSkip(context.Background(), slog.LevelInfo, 2, msg, args...)
}

func (s *SlogAdapter) Warn(msg string, args ...any) {
	s.logWithSkip(context.Background(), slog.LevelWarn, 2, msg, args...)
}

func (s *SlogAdapter) Error(msg string, args ...any) {
	s.logWithSkip(context.Background(), slog.LevelError, 2, msg, args...)
}

func (s *SlogAdapter) DebugCtx(ctx context.Context, msg string, args ...any) {
	s.buildContextualLogger(ctx, slog.LevelDebug, msg, args...)
}

func (s *SlogAdapter) InfoCtx(ctx context.Context, msg string, args ...any) {
	s.buildContextualLogger(ctx, slog.LevelInfo, msg, args...)
}

func (s *SlogAdapter) WarnCtx(ctx context.Context, msg string, args ...any) {
	s.buildContextualLogger(ctx, slog.LevelWarn, msg, args...)
}

func (s *SlogAdapter) ErrorCtx(ctx context.Context, msg string, args ...any) {
	s.buildContextualLogger(ctx, slog.LevelError, msg, args...)
}

func (s *SlogAdapter) With(args ...any) Logger {
	return &SlogAdapter{l: s.l.With(args...)}
}

func (s *SlogAdapter) logWithSkip(
	ctx context.Context,
	level slog.Level,
	skip int,
	msg string,
	args ...any,
) {
	if !s.l.Enabled(ctx, level) {
		return
	}

	// Get caller information
	var pcs [1]uintptr
	// skip 2 to capture the correct caller frame:
	// Frame 1: runtime.Callers (this call)
	// Frame 2: logWithSkip (this method)
	runtime.Callers(skip+2, pcs[:])

	// Create a new record with proper source
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	// Call handler with the record - this preserves With/WithGroup
	_ = s.l.Handler().Handle(ctx, r)
}

func (s *SlogAdapter) buildContextualLogger(
	ctx context.Context,
	level slog.Level,
	msg string,
	args ...any,
) {
	if l := GetFromContext(ctx); l != nil {
		if sa, ok := l.(*SlogAdapter); ok {
			// addtional 3 skips to capture the correct caller frame:
			// Frame 3: this function
			// Frame 4: wrapper (InfoCtx, DebugCtx ...)
			// Frame 5: Actual Caller
			sa.logWithSkip(ctx, level, 3, msg, args...)
			return
		}
	}
	s.logWithSkip(ctx, level, 3, msg, args...)
}
