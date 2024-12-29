package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Upstairs struct {
	LuxSensor      *hal.LightSensor
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
}

func newUpstairs() Upstairs {
	return Upstairs{
		LuxSensor:      hal.NewLightSensor("sensor.presence_sensor_fp2_b6d8_light_sensor_light_level"),
		PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_1"),
	}
}

func (u *Upstairs) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Upstairs bookshelf lamps").
			WithSensors(home.Upstairs.PresenceSensor).
			WithLights(
				home.LivingRoom.MoroccanLamp,
				home.LivingRoom.SaltLamp,
			).
			WithBrightness(64).
			TurnsOffAfter(15 * time.Minute),
	}
}
