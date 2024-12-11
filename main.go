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
	Study   Study
}

func (m *Marnixkade) TestFunc() {
	slog.Info("TestFunc")
}

type Hallway struct {
	Lights       *hal.Light
	MotionSensor *hal.Entity
}

type Study struct {
	Lights         *hal.Light  // Main lights
	LuxSensor      *hal.Entity // Aqara FP2
	PresenceSensor *hal.Entity // Aqara FP2
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
		Study: Study{
			Lights:         hal.NewLight("light.study_lights"),
			PresenceSensor: hal.NewEntity("binary_sensor.presence_sensor_fp2_11ad_presence_sensor_1"),
			LuxSensor:      hal.NewEntity("sensor.presence_sensor_fp2_11ad_light_sensor_light_level"),
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
			TurnsOffAfter(5*time.Minute),

		halautomations.NewSensorsTriggersLights().
			WithName("Study lights").
			WithSensors(home.Study.PresenceSensor).
			WithLights(home.Study.Lights).
			TurnsOffAfter(5*time.Minute),

		// hal.NewAutomation().
		// 	WithName("Hallway lights").
		// 	WithEntities(home.Hallway.MotionSensor).
		// 	WithAction(home.TestFunc),
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
