package main

import (
	"context"
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type DiningRoom struct {
	Lights         *hal.Light
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
}

func newDiningRoom() DiningRoom {
	return DiningRoom{
		Lights:         hal.NewLight("light.dining"),
		PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_2"),
	}
}

func (d *DiningRoom) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Dining table lights").
			WithCondition(func() bool {
				return !home.NightMode.IsOn() // Don't auto turn on lights if night mode is on
			}).
			WithSensors(home.DiningRoom.PresenceSensor).
			WithLights(home.DiningRoom.Lights).
			SetScene(brightLight).
			TurnsOffAfter(15 * time.Minute),
		// WithHumanOverrideFor(6 * 60 * time.Minute),

		hal.NewAutomation().
			WithName("Night mode dining lamp").
			WithEntities(home.DiningRoom.PresenceSensor).
			WithAction(func(ctx context.Context, _ hal.EntityInterface) {
				if !home.NightMode.IsOn() {
					return
				}

				// At night only the Moroccan lamp comes on when someone is at
				// the dining table, and it turns off as soon as nobody is
				// there anymore.
				if home.DiningRoom.PresenceSensor.IsOn() {
					home.LivingRoom.MoroccanLamp.TurnOnContext(ctx, map[string]any{"brightness": 64})
				} else {
					home.LivingRoom.MoroccanLamp.TurnOffContext(ctx)
				}
			}),
	}
}
