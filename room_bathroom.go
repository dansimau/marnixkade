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

	SwitchOffButton *hal.Button
}

func newBathroom() Bathroom {
	return Bathroom{
		Fan:          hal.NewLight("light.bathroom_fan"),
		MotionSensor: hal.NewBinarySensor("binary_sensor.bathroom_sensor_motion"),
		Light:        hal.NewLight("light.bathroom"),

		SwitchOffButton: hal.NewButton("event.bathroom_switch_button_4"),
	}
}

func (room *Bathroom) LightIsOff() bool {
	return !room.Light.IsOn()
}

func (room *Bathroom) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		hal.NewAutomation().
			WithName("Button switch off bathroom fan").
			WithEntities(room.SwitchOffButton).
			WithAction(func(trigger hal.EntityInterface) {
				if room.SwitchOffButton.PressedTimes() > 1 {
					room.Fan.TurnOff()
				}
			}),

		halautomations.NewSensorsTriggerLights().
			WithName("Bathroom light").
			WithSensors(room.MotionSensor).
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, nightLight).
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			WithLights(room.Light).
			TurnsOffAfter(15 * time.Minute),
		// WithHumanOverrideFor(40 * time.Minute),

		// Turn on bathroom fan 1 minute after lights go on (i.e. if someone is
		// lingering in the bathroom)
		halautomations.NewTimer("Bathroom fan on timer").
			WithEntities(room.Light).
			Condition(room.Light.IsOn).
			Condition(home.NightMode.IsOff). // Don't turn fan on at night because it is noisy
			Duration(1 * time.Minute).
			Run(func() {
				room.Fan.TurnOn()
			}),

		// Turn bathroom fan off 40 mins after lights go off
		halautomations.NewTimer("Bathroom fan off timer").
			WithEntities(room.Light).
			Condition(room.LightIsOff).
			Duration(40 * time.Minute).
			Run(func() {
				room.Fan.TurnOff()
			}),
	}
}
