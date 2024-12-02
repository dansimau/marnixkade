package hal

import (
	"errors"

	"github.com/dansimau/home-automation/pkg/hassws"
)

type Light struct {
	*Entity
}

func NewLight(id string) *Light {
	return &Light{Entity: NewEntity(id)}
}

func (c *Light) IsOn() bool {
	return c.Entity.GetState().State == "on"
}

func (l *Light) TurnOn() error {
	if l.connection == nil {
		return errors.New("entity not registered")
	}

	_, err := l.connection.Client.CallService(hassws.CallServiceRequest{
		Type:    hassws.MessageTypeCallService,
		Domain:  "light",
		Service: "turn_on",
		Data: map[string]any{
			"entity_id": []string{l.GetID()},
		},
	})

	return err
}

func (l *Light) TurnOff() error {
	if l.connection == nil {
		return errors.New("entity not registered")
	}

	_, err := l.connection.Client.CallService(hassws.CallServiceRequest{
		Type:    hassws.MessageTypeCallService,
		Domain:  "light",
		Service: "turn_off",
		Data: map[string]any{
			"entity_id": []string{l.GetID()},
		},
	})

	return err
}
