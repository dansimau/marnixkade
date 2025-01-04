package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Bathroom struct {
	Fan          *hal.Light
	MotionSensor *hal.BinarySensor
	Light        *hal.Light
}

func newBathroom() Bathroom {
	return Bathroom{
		Fan:          hal.NewLight("light.bathroom_fan"),
		MotionSensor: hal.NewBinarySensor("binary_sensor.bathroom_motion"),
		Light:        hal.NewLight("light.bathroom"),
	}
}

func (room *Bathroom) LightIsOff() bool {
	return !room.Light.IsOn()
}

func (room *Bathroom) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Bathroom light").
			WithSensors(room.MotionSensor).
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, nightLight).
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			WithLights(room.Light).
			TurnsOffAfter(15 * time.Minute).
			WithHumanOverrideFor(40 * time.Minute),

		// Turn on bathroom fan 1 minute after lights go on (i.e. if someone is
		// lingering in the bathroom)
		halautomations.NewTimer("Bathroom fan on timer").
			WithEntities(room.Fan).
			Condition(room.Light.IsOn).
			Duration(1 * time.Minute).
			Run(func() {
				room.Fan.TurnOn()
			}),

		// Turn bathroom fan off 40 mins after lights go off
		halautomations.NewTimer("Bathroom fan off timer").
			WithEntities(room.Fan).
			Condition(room.LightIsOff).
			Duration(40 * time.Minute).
			Run(func() {
				room.Fan.TurnOff()
			}),
	}
}
