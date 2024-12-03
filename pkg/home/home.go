package home

import (
	"time"

	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/halautomations"
	"github.com/dansimau/home-automation/pkg/hassws"
)

type Marnixkade struct {
	*hal.Connection

	Hallway Hallway
}

func NewMarnixkade() *Marnixkade {
	home := &Marnixkade{
		Connection: hal.NewConnection(hassws.NewWebsocketAPI(HomeAssistantConfig)),
		Hallway:    hallway,
	}

	// Walk the struct and find/register all entities
	home.FindEntities(home)

	// Register automations
	home.RegisterAutomations(
		halautomations.NewSensorsTriggersLights().
			WithName("Hallway lights").
			WithSensors(home.Hallway.MotionSensor).
			WithLights(home.Hallway.Lights).
			TurnsOffAfter(5 * time.Minute),
	)

	return home
}
