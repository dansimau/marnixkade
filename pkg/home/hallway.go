package home

import "github.com/dansimau/home-automation/pkg/hal"

type Hallway struct {
	Lights         *hal.Light
	MotionPresence *hal.MotionSensor
}
