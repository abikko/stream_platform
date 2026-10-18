package main

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
}
