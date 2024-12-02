package home

import (
	"time"

	"github.com/dansimau/home-automation/pkg/automations"
	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/hassws"
)

type Marnixkade struct {
	*hal.Connection

	Hallway Hallway
}

func NewHome() *Marnixkade {
	// Construct your home
	home := &Marnixkade{
		Connection: hal.NewConnection(hassws.NewWebsocketAPI(HomeAssistantConfig)),

		Hallway: Hallway{
			Lights:       hal.NewLight("light.front_hallway"),
			MotionSensor: hal.NewEntity("binary_sensor.hallway_motion"),
		},
	}

	// Register entities (TODO: Init function can walk the struct and do this automatically with reflection)
	home.RegisterEntities(
		home.Hallway.MotionSensor,
		home.Hallway.Lights,
	)

	// Register automations
	home.RegisterAutomations(
		(&automations.SensorsTriggersLights{
			Sensors:       hal.Entities{home.Hallway.MotionSensor},
			Lights:        []*hal.Light{home.Hallway.Lights},
			TurnsOffAfter: 5 * time.Minute,
		}).Build(),
	)

	return home
}
