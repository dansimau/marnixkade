package hal_test

import (
	"testing"

	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/haltest"
	"github.com/dansimau/home-automation/pkg/homeassistant"
)

func TestConnection(t *testing.T) {
	conn, server, cleanup := haltest.NewClientServer(t)
	defer cleanup()

	// Create test entity and register it
	entity := hal.NewEntity("test.entity")
	conn.RegisterEntities(entity)

	// Send state change event
	server.SendStateChangeEvent(homeassistant.Event{
		EventData: homeassistant.EventData{
			EntityID: "test.entity",
			NewState: &homeassistant.State{State: "on"},
		},
	})

	// Verify entity state was updated
	haltest.WaitFor(t, func() bool {
		return entity.State.State == "on"
	})
}
