package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Downstairs struct {
	AllLights *hal.Light

	MotionSensorStairs *hal.BinarySensor
	MotionSensorWindow *hal.BinarySensor
}

func newDownstairs() Downstairs {
	return Downstairs{
		AllLights: hal.NewLight("light.downstairs"),

		MotionSensorStairs: hal.NewBinarySensor("binary_sensor.stairs_sensor_motion"),
		MotionSensorWindow: hal.NewBinarySensor("binary_sensor.downstairs_sensor_motion"),
	}
}

func (d *Downstairs) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Downstairs lights").
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, nightLight).
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			WithSensors(
				home.Downstairs.MotionSensorStairs,
				home.Downstairs.MotionSensorWindow,
			).
			WithLights(home.Downstairs.AllLights).
			TurnsOffAfter(5 * time.Minute),
	}
}
