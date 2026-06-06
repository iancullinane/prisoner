package logging

import (
	"io"
	"log/slog"
	"os"
)

func New(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func NewLogger() *slog.Logger {
	return New(os.Stdout)
}
