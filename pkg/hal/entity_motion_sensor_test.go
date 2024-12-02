package hal_test

import (
	"testing"
	"time"

	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/homeassistant"
	"github.com/dansimau/home-automation/pkg/testutil"
)

func TestPhysicalMotionSensorTriggered(t *testing.T) {
	t.Parallel()

	conn, server, cleanup := newClientServer(t)
	defer cleanup()

	// Create test entity and register it
	motionSensor := hal.NewMotionSensor("test.motion_sensor", hal.MotionSensorConfig{
		ResetTimeout: 1 * time.Second,
	})
	conn.RegisterEntities(motionSensor)

	// Send state change event
	server.SendStateChangeEvent(homeassistant.Event{
		EventData: homeassistant.EventData{
			EntityID: "test.motion_sensor",
			NewState: &homeassistant.State{State: "on"},
		},
	})

	// Verify entity state was updated
	testutil.WaitFor(t, func() bool {
		return motionSensor.State.State == "on"
	})
}

func TestPhysicalMotionSensorResetIgnored(t *testing.T) {
	t.Parallel()

	resetTimeout := 3 * time.Second

	conn, server, cleanup := newClientServer(t)
	defer cleanup()

	// Create test entity and register it
	motionSensor := hal.NewMotionSensor("test.motion_sensor", hal.MotionSensorConfig{
		ResetTimeout: resetTimeout,
	})
	conn.RegisterEntities(motionSensor)

	// Motion sensor is triggered
	server.SendStateChangeEvent(homeassistant.Event{
		EventData: homeassistant.EventData{
			EntityID: "test.motion_sensor",
			NewState: &homeassistant.State{State: "on"},
		},
	})

	// Verify entity state was updated
	testutil.WaitFor(t, func() bool {
		return motionSensor.Triggered() == true
	}, "motion sensor should be triggered")

	// Send motion sensor reset
	server.SendStateChangeEvent(homeassistant.Event{
		EventData: homeassistant.EventData{
			EntityID: "test.motion_sensor",
			NewState: &homeassistant.State{State: "off"},
		},
	})

	// Motion sensor should still be triggered because reset timeout has not elapsed
	testutil.WaitFor(t, func() bool {
		return motionSensor.Triggered() == true
	}, "motion sensor should still be triggered")

	// Wait for reset timeout to elapse
	time.Sleep(resetTimeout)

	// Motion sensor should no longer be triggered
	testutil.WaitFor(t, func() bool {
		return motionSensor.Triggered() == false
	}, "motion sensor should no longer be triggered")
}
