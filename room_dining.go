package main

import (
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
			WithSensors(home.DiningRoom.PresenceSensor).
			WithLights(home.DiningRoom.Lights).
			SetScene(brightLight).
			TurnsOffAfter(15 * time.Minute).
			WithHumanOverrideFor(6 * 60 * time.Minute),
	}
}
