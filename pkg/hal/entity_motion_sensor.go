package hal

import (
	"time"

	"github.com/dansimau/home-automation/pkg/hassws"
	"github.com/dansimau/home-automation/pkg/homeassistant"
)

type MotionSensor struct {
	*Entity

	cfg MotionSensorConfig

	timer     *Timer
	triggered bool
}

type MotionSensorConfig struct {
	ResetTimeout time.Duration
}

func NewMotionSensor(id string, cfg MotionSensorConfig) *MotionSensor {
	m := &MotionSensor{
		Entity: NewEntity(id),
		cfg:    cfg,
	}
	m.timer = NewTimer(m.resetMotionTriggered)
	return m
}

func (s *MotionSensor) resetMotionTriggered() {
	s.triggered = false

	state := s.Entity.GetState()
	state.State = "off"

	s.connection.StateChangeEvent(hassws.EventMessage{
		Event: homeassistant.Event{
			EventData: homeassistant.EventData{
				EntityID: s.GetID(),
				NewState: &state,
			},
		},
	})
}

func (s *MotionSensor) Triggered() bool {
	return s.triggered
}

func (s *MotionSensor) SetState(state homeassistant.State) {
	if state.State == "on" {
		s.triggered = true
		s.timer.Reset(s.cfg.ResetTimeout)
	}
}
