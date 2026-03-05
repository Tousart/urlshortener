package app

import (
	"log/slog"
	"os"
)

// InitLogger() - инициализация логгера
func InitLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
