package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Hallway struct {
	Lights       *hal.Light
	MotionSensor *hal.Entity
}

func newHallway() Hallway {
	return Hallway{
		Lights:       hal.NewLight("light.front_hallway"),
		MotionSensor: hal.NewEntity("binary_sensor.hallway_motion"),
	}
}

func (h *Hallway) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Front hallway lights").
			WithSensors(home.Hallway.MotionSensor).
			WithLights(home.Hallway.Lights).
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, map[string]any{"brightness": 50}). // 20% brightness at night
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			TurnsOffAfter(1 * time.Minute),
	}
}
