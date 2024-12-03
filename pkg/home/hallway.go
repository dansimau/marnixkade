package home

import "github.com/dansimau/home-automation/pkg/hal"

type Hallway struct {
	Lights       *hal.Light
	MotionSensor *hal.Entity
}

var hallway = Hallway{
	Lights:       hal.NewLight("light.front_hallway"),
	MotionSensor: hal.NewEntity("binary_sensor.hallway_motion"),
}
