package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init() {
	Logger = slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}
