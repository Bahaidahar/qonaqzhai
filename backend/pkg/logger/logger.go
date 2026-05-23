// Package logger provides a slog setup that every service shares.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a structured slog.Logger tagged with the service name. The level
// is read from LOG_LEVEL (debug/info/warn/error); default is info.
func New(service string) *slog.Logger {
	lvl := parseLevel(os.Getenv("LOG_LEVEL"))
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(h).With(slog.String("svc", service))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
