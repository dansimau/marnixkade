package main

import (
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.AnyValue(time.Now().Format("2006-01-02 15:04:05.000")),
				}
			}
			return a
		},
	})))

	if err := NewMarnixkade().Start(); err != nil {
		slog.Error("Error starting home", "error", err)
		os.Exit(1)
	}

	select {}
}
