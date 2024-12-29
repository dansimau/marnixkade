package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	if err := NewMarnixkade().Start(); err != nil {
		slog.Error("Error starting home", "error", err)
		os.Exit(1)
	}

	select {}
}
