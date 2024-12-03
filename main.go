package main

import (
	"log/slog"
	"os"

	"github.com/dansimau/home-automation/pkg/home"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	if err := home.NewMarnixkade().Start(); err != nil {
		slog.Error("Error starting home", "error", err)
		os.Exit(1)
	}

	select {}
}
