// Package logger configures the application's slog logger.
package logger

import (
	"log/slog"
	"os"
)

// New returns a JSON slog handler at the given level.
// Empty level defaults to "info".
func New(level string) *slog.Logger {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(h)
}
