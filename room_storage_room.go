package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type StorageRoom struct {
	Lights       *hal.Light
	MotionSensor *hal.BinarySensor
}

func newStorageRoom() StorageRoom {
	return StorageRoom{
		Lights:       hal.NewLight("light.storage_room"),
		MotionSensor: hal.NewBinarySensor("binary_sensor.hue_motion_sensor_1_motion"),
	}
}

func (room *StorageRoom) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Storage room lights").
			WithSensors(room.MotionSensor).
			WithLights(room.Lights).
			TurnsOffAfter(5 * time.Minute),
	}
}
