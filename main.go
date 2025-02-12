package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level: slog.LevelDebug,
		}),
	))

	if err := NewMarnixkade().Start(); err != nil {
		slog.Error("Error starting home", "error", err)
		os.Exit(1)
	}

	select {}
}
