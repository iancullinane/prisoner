package logging

import (
	"io"
	"log/slog"
	"strings"
)

// New builds a slog.Logger writing to w, filtered at level, using the given
// format. format is "text" (human-readable, the dev default) or "json"
// (structured); anything else falls back to text.
func New(w io.Writer, level slog.Level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	switch strings.ToLower(format) {
	case "json":
		h = slog.NewJSONHandler(w, opts)
	default:
		h = slog.NewTextHandler(w, opts)
	}

	return slog.New(h)
}

// ParseLevel maps a config/flag string to a slog.Level. Unknown values default
// to debug because the application runs verbose in development.
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
