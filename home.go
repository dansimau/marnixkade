package main

import (
	"log/slog"
	"os"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Marnixkade struct {
	*hal.Connection

	DiningRoom DiningRoom
	Downstairs Downstairs
	Bedroom    Bedroom
	Hallway    Hallway
	Kitchen    Kitchen
	LivingRoom LivingRoom
	Study      Study
	Upstairs   Upstairs

	// Guest mode is a switch that can be used to turn off certain automations
	// when guests are over.
	GuestMode *hal.InputBoolean

	// NightMode is the "bed time" switch that controls lights downstairs and
	// in the bedrooms.
	NightMode *hal.InputBoolean
}

func NewMarnixkade() *Marnixkade {
	cfg, err := hal.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	home := &Marnixkade{
		Connection: hal.NewConnection(*cfg),

		Bedroom:    newBedroom(),
		DiningRoom: newDiningRoom(),
		Downstairs: newDownstairs(),
		Hallway:    newHallway(),
		Kitchen:    newKitchen(),
		LivingRoom: newLivingRoom(),
		Study:      newStudy(),
		Upstairs:   newUpstairs(),

		GuestMode: hal.NewInputBoolean("input_boolean.guest_mode"),
		NightMode: hal.NewInputBoolean("input_boolean.bedtime_switch"),
	}

	// Walk the struct and find/register all entities
	home.FindEntities(home)

	// Register automations
	home.RegisterAutomations(home.Bedroom.Automations(home)...)
	home.RegisterAutomations(home.DiningRoom.Automations(home)...)
	home.RegisterAutomations(home.Downstairs.Automations(home)...)
	home.RegisterAutomations(home.Hallway.Automations(home)...)
	home.RegisterAutomations(home.LivingRoom.Automations(home)...)
	home.RegisterAutomations(home.Study.Automations(home)...)
	home.RegisterAutomations(home.Upstairs.Automations(home)...)

	home.RegisterAutomations(halautomations.NewPrintDebug("Debug", hal.NewEntity("event.main_switch_button_1")))

	return home
}
