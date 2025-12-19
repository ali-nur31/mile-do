package logger

import (
	"log/slog"
	"os"
)

func InitializeLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
