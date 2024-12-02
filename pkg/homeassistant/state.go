package homeassistant

import (
	"encoding/json"
	"time"
)

const (
	EventTypeStateChanged = "state_changed"
)

type State struct {
	EntityID string `json:"entity_id"`

	State      string `json:"state"`
	Attributes struct {
		DeviceClass  string `json:"device_class"`
		FriendlyName string `json:"friendly_name"`

		json.RawMessage
	} `json:"attributes"`

	LastChanged  time.Time `json:"last_changed"`
	LastReported time.Time `json:"last_reported"`
	LastUpdated  time.Time `json:"last_updated"`
}

// func MergeState(state *State, newState State) {
// 	if newState.EntityID != "" {
// 		state.EntityID = newState.EntityID
// 	}

// 	if newState.State != "" {
// 		state.State = newState.State
// 	}

// 	if newState.Attributes.StateClass != "" {
// 		state.Attributes.StateClass = newState.Attributes.StateClass
// 	}
// 	if newState.Attributes.UnitOfMeasurement != "" {
// 		state.Attributes.UnitOfMeasurement = newState.Attributes.UnitOfMeasurement
// 	}
// 	if newState.Attributes.DeviceClass != "" {
// 		state.Attributes.DeviceClass = newState.Attributes.DeviceClass
// 	}
// 	if newState.Attributes.FriendlyName != "" {
// 		state.Attributes.FriendlyName = newState.Attributes.FriendlyName
// 	}
// }
