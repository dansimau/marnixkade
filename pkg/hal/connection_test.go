package hal_test

import (
	"testing"

	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/homeassistant"
	"github.com/dansimau/home-automation/pkg/testutil"
)

func TestConnection(t *testing.T) {
	conn, server, cleanup := newClientServer(t)
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
	testutil.WaitFor(t, func() bool {
		return entity.State.State == "on"
	})
}
