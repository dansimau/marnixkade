package hal

import (
	"github.com/dansimau/home-automation/pkg/homeassistant"
)

// EntityLike is the interface that we accept which allows us to create
// custom components that embed an Entity.
type EntityLike interface {
	ConnectionBinder

	GetID() string
	GetState() homeassistant.State
	SetState(event homeassistant.State)
}

type Entities []EntityLike

// Entity is a base type for all entities that can be embedded into other types.
type Entity struct {
	connection *Connection

	homeassistant.State
}

func NewEntity(id string) *Entity {
	return &Entity{State: homeassistant.State{EntityID: id}}
}

// BindConnection binds the entity to the connection. This allows entities to
// publish messages.
func (e *Entity) BindConnection(connection *Connection) {
	e.connection = connection
}

func (e *Entity) GetID() string {
	return e.State.EntityID
}

func (e *Entity) SetState(state homeassistant.State) {
	e.State = state
}

func (e *Entity) GetState() homeassistant.State {
	return e.State
}
