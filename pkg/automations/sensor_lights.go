package automations

import (
	"time"

	"github.com/dansimau/home-automation/pkg/hal"
)

type SensorsTriggersLights struct {
	Sensors       []hal.EntityLike
	Lights        []*hal.Light
	TurnsOffAfter time.Duration

	turnOffTimer *time.Timer
}

// triggered returns true if any of the sensors have been triggered
func (a *SensorsTriggersLights) triggered() bool {
	for _, sensor := range a.Sensors {
		if sensor.GetState().State == "on" {
			return true
		}
	}
	return false
}

func (a *SensorsTriggersLights) action() {
	if a.triggered() {
		a.stopTurnOffTimer()
		a.turnOnLights()
	} else {
		a.startTurnOffTimer()
	}
}

func (a *SensorsTriggersLights) startTurnOffTimer() {
	if a.turnOffTimer == nil {
		a.turnOffTimer = time.AfterFunc(a.TurnsOffAfter, a.turnOffLights)
	} else {
		a.turnOffTimer.Reset(a.TurnsOffAfter)
	}
}

// stop timer
func (a *SensorsTriggersLights) stopTurnOffTimer() {
	if a.turnOffTimer != nil {
		a.turnOffTimer.Stop()
	}
}

// turn on lights
func (a *SensorsTriggersLights) turnOnLights() {
	for _, light := range a.Lights {
		light.TurnOn()
	}
}

// turn off lights
func (a *SensorsTriggersLights) turnOffLights() {
	for _, light := range a.Lights {
		light.TurnOff()
	}
}

func (a *SensorsTriggersLights) Build() hal.Automation {
	return hal.Automation{
		Name:     "Motion triggers lights",
		Entities: hal.Entities(a.Sensors),
		Action:   a.action,
	}
}
