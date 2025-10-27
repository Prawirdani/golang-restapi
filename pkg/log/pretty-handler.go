package log

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
)

// PrettyHandler wraps slog.Handler to provide colorful, formatted output
type PrettyHandler struct {
	handler slog.Handler
	attrs   []slog.Attr
	groups  []string
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorCyan   = "\033[36m"
)

// formatValue formats slog.Value, handling byte slices specially
func formatValue(v slog.Value) string {
	// Check if it's a byte slice
	if v.Kind() == slog.KindAny {
		if b, ok := v.Any().([]byte); ok {
			// Try to convert to string if it's printable
			if len(b) > 0 && isPrintable(b) {
				return string(b)
			}
			// Otherwise show hex for binary data
			return "0x" + string(b[:min(len(b), 32)]) + "..."
		}

		// Try to marshal complex types to JSON for readability
		val := v.Any()
		if val != nil {
			// Check if it's a struct or complex type
			if jsonBytes, err := json.Marshal(val); err == nil {
				// Return the JSON string if it's valid and printable
				if isPrintable(jsonBytes) {
					return string(jsonBytes)
				}
			}
		}
	}
	return v.String()
}

// isPrintable checks if byte slice contains printable characters (like JSON)
func isPrintable(b []byte) bool {
	for _, c := range b {
		if c < 32 && c != '\n' && c != '\r' && c != '\t' {
			return false
		}
		if c > 126 {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	// Color code the level
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = colorGray
	case slog.LevelInfo:
		levelColor = colorBlue
	case slog.LevelWarn:
		levelColor = colorYellow
	case slog.LevelError:
		levelColor = colorRed
	}

	// Format: TIME LEVEL MESSAGE key=value key=value
	timeStr := r.Time.Format("15:04:05.000")

	// Print formatted log
	os.Stdout.WriteString(colorGray + timeStr + colorReset + " ")
	os.Stdout.WriteString(levelColor + level + colorReset + " ")
	os.Stdout.WriteString(r.Message)

	// Print handler's persistent attributes (without group prefix - they were added before groups)
	for _, field := range h.attrs {
		os.Stdout.WriteString(" " + colorCyan + field.Key + colorReset + "=")
		os.Stdout.WriteString(formatValue(field.Value))
	}

	// Print record's attributes (with group prefix - they're added in the context of the group)
	r.Attrs(func(a slog.Attr) bool {
		key := a.Key
		// Add group prefix if any
		if len(h.groups) > 0 {
			prefix := ""
			for _, g := range h.groups {
				if prefix != "" {
					prefix += "."
				}
				prefix += g
			}
			key = prefix + "." + key
		}
		os.Stdout.WriteString(" " + colorCyan + key + colorReset + "=")
		os.Stdout.WriteString(formatValue(a.Value))
		return true
	})

	os.Stdout.WriteString("\n")

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create new handler with accumulated attributes
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &PrettyHandler{
		handler: h.handler.WithAttrs(attrs),
		attrs:   newAttrs,
		groups:  h.groups,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// Create new handler with accumulated groups
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &PrettyHandler{
		handler: h.handler.WithGroup(name),
		attrs:   h.attrs,
		groups:  newGroups,
	}
}

// NewPrettyHandler creates a new pretty handler for development
func NewPrettyHandler(opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	return &PrettyHandler{
		handler: slog.NewTextHandler(os.Stdout, opts),
		attrs:   []slog.Attr{},
		groups:  []string{},
	}
}
