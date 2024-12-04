package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
	"github.com/dansimau/hal/hassws"
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
	home := &Marnixkade{
		Connection: hal.NewConnection(hassws.NewWebsocketAPI(hassws.ClientConfig{
			Host:  "nas:8123",
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIyMmI4YmMwNWQ2YTM0YmNiYTYyMmE3MGI2M2Y4ZWU2NiIsImlhdCI6MTczMjc0Njc4OSwiZXhwIjoyMDQ4MTA2Nzg5fQ.xi8DDzg50D-if-I0j4q-r-TzQ__xVl-13tcB5_hUmRQ",
		})),
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
