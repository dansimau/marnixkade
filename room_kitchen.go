package main

import "github.com/dansimau/hal"

type Kitchen struct {
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
}

func newKitchen() Kitchen {
	return Kitchen{
		PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_4"),
	}
}
