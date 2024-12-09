package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Marnixkade struct {
	*hal.Connection

	Hallway Hallway
}

type Hallway struct {
	Lights       *hal.Light
	MotionSensor *hal.Entity
}

func NewMarnixkade() *Marnixkade {
	cfg, err := hal.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	home := &Marnixkade{
		Connection: hal.NewConnection(*cfg),
		Hallway: Hallway{
			Lights:       hal.NewLight("light.front_hallway"),
			MotionSensor: hal.NewEntity("binary_sensor.hallway_motion"),
		},
	}

	// Walk the struct and find/register all entities
	home.FindEntities(home)

	// Register automations
	home.RegisterAutomations(
		halautomations.NewSensorsTriggersLights().
			WithName("Hallway lights").
			WithSensors(home.Hallway.MotionSensor).
			WithLights(home.Hallway.Lights).
			TurnsOffAfter(5 * time.Minute),
	)

	return home
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	if err := NewMarnixkade().Start(); err != nil {
		slog.Error("Error starting home", "error", err)
		os.Exit(1)
	}

	select {}
}
