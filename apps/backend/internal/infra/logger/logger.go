package logger

import (
	"log/slog"
	"os"
)

func New(appName string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})).With("app", appName)
}
