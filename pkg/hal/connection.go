package hal

import (
	"log/slog"
	"sync"

	"github.com/dansimau/home-automation/pkg/hassws"
	"github.com/davecgh/go-spew/spew"
)

type Connection struct {
	*hassws.Client

	automations map[string][]Automation
	entities    map[string]EntityLike

	mutex sync.RWMutex
}

type ConnectionBinder interface {
	BindConnection(connection *Connection)
}

func NewConnection(api *hassws.Client) *Connection {
	return &Connection{
		Client: api,

		automations: make(map[string][]Automation),
		entities:    make(map[string]EntityLike),
	}
}

func (h *Connection) RegisterAutomations(automations ...Automation) {
	for _, automation := range automations {
		for _, entity := range automation.Entities {
			h.automations[entity.GetID()] = append(h.automations[entity.GetID()], automation)
		}
	}
}

func (h *Connection) RegisterEntities(entities ...EntityLike) {
	for _, entity := range entities {
		entity.BindConnection(h)
		h.entities[entity.GetID()] = entity
	}
}

func (h *Connection) Start() error {
	if err := h.Client.Connect(); err != nil {
		return err
	}

	return h.Client.SubscribeEvents(string(hassws.MessageTypeStateChanged), h.StateChangeEvent)
}

// Process incoming state change events. Dispatch state change to the relevant
// entity and fire any automations listening for state changes to this entity.
func (h *Connection) StateChangeEvent(event hassws.EventMessage) {
	// defer perf.Timer(func(timeTaken time.Duration) {
	// 	slog.Debug("Tick processing time (total)", "duration", timeTaken)
	// })()

	h.mutex.Lock()
	defer h.mutex.Unlock()

	// defer perf.Timer(func(timeTaken time.Duration) {
	// 	slog.Debug("Tick processing time (processing)", "duration", timeTaken)
	// })()

	entity, ok := h.entities[event.Event.EventData.EntityID]
	if !ok {
		slog.Debug("Entity not registered", "EntityID", event.Event.EventData.EntityID)
		return
	}

	slog.Debug("State changed for", "EntityID", event.Event.EventData.EntityID, "NewState", spew.Sdump(event.Event.EventData.NewState))

	if event.Event.EventData.NewState != nil {
		entity.SetState(*event.Event.EventData.NewState)
	}

	// Dispatch automations
	for _, automation := range h.automations[event.Event.EventData.EntityID] {
		automation.Action()
	}
}
