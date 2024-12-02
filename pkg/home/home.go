package home

import (
	"time"

	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/hassws"
)

type Marnixkade struct {
	*hal.Connection

	Hallway Hallway
}

var Home *Marnixkade = &Marnixkade{
	Hallway: Hallway{
		Lights: hal.NewLight("light.front_hallway"),
		MotionPresence: hal.NewMotionSensor("binary_sensor.hallway_motion", hal.MotionSensorConfig{
			ResetTimeout: 5 * time.Minute,
		}),
	},
}

func NewHome() *Marnixkade {
	// Construct your home
	home := &Marnixkade{
		Connection: hal.NewConnection(hassws.NewWebsocketAPI(HomeAssistantConfig)),

		Hallway: Hallway{
			Lights: hal.NewLight("light.front_hallway"),
			MotionPresence: hal.NewMotionSensor("binary_sensor.hallway_motion", hal.MotionSensorConfig{
				ResetTimeout: time.Second * 5,
			}),
		},
	}

	// Register entities (TODO: Init function can walk the struct and do this automatically with reflection)
	home.RegisterEntities(
		home.Hallway.MotionPresence,
		home.Hallway.Lights,
	)

	// Register automations
	home.RegisterAutomations(
		hal.Automation{
			Name:     "Hallway: Motion (Day)",
			Entities: hal.Entities{home.Hallway.MotionPresence},
			Action: func() {
				if home.Hallway.MotionPresence.Triggered() {
					home.Hallway.Lights.TurnOn()
				} else {
					home.Hallway.Lights.TurnOff()
				}
			},
		},
	)

	return home
}
