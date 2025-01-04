package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Kitchen struct {
	MotionSensor   *hal.BinarySensor
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
	StripLight     *hal.Light
}

func newKitchen() Kitchen {
	return Kitchen{
		MotionSensor:   hal.NewBinarySensor("binary_sensor.kitchen_motion"),
		PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_4"),
		StripLight:     hal.NewLight("light.kitchen_strip"),
	}
}

func (k *Kitchen) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Kitchen strip light").
			WithSensors(k.MotionSensor).
			WithLights(k.StripLight).
			TurnsOffAfter(15 * time.Minute),
	}
}
