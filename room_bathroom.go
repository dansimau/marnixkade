package main

import (
	"context"
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Bathroom struct {
	Fan          *hal.Light
	MotionSensor *hal.BinarySensor
	MotionAware  *hal.BinarySensor
	Light        *hal.Light

	SwitchOffButton *hal.Button
}

func newBathroom() Bathroom {
	return Bathroom{
		Fan:          hal.NewLight("light.bathroom_fan"),
		MotionSensor: hal.NewBinarySensor("binary_sensor.bathroom_sensor_motion"),
		MotionAware:  hal.NewBinarySensor("binary_sensor.bathroom_motionaware"),
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
			WithAction(func(ctx context.Context, trigger hal.EntityInterface) {
				if room.SwitchOffButton.PressedTimes() > 1 {
					room.Fan.TurnOffContext(ctx)
				}
			}),

		halautomations.NewSensorsTriggerLights().
			WithName("Bathroom light").
			WithSensors(
				home.Bathroom.MotionSensor,
				home.Bathroom.MotionAware,
			).
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, nightLight).
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			// WithConditionScene(func() bool { return true }, spookyLight).
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
			Run(func(ctx context.Context) {
				room.Fan.TurnOnContext(ctx)
			}),

		// Turn bathroom fan off automatically after a while
		halautomations.NewTimer("Bathroom fan off timer").
			WithEntities(room.Light).
			Condition(room.LightIsOff).
			Duration(90 * time.Minute).
			Run(func(ctx context.Context) {
				room.Fan.TurnOffContext(ctx)
			}),
	}
}
