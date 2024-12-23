package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Marnixkade struct {
	*hal.Connection

	DiningRoom DiningRoom
	Hallway    Hallway
	Kitchen    Kitchen
	LivingRoom LivingRoom
	Study      Study
	Upstairs   Upstairs

	// Guest mode is a switch that can be used to turn off certain automations
	// when guests are over.
	GuestMode *hal.BinarySensor
}

type DiningRoom struct {
	Lights         *hal.Light
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
}

type Hallway struct {
	Lights       *hal.Light
	MotionSensor *hal.Entity
}

type Kitchen struct {
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
}

type LivingRoom struct {
	ArcherLamp   *hal.Light
	MainLights   hal.LightGroup // Hue Filament Bulbs
	MoroccanLamp *hal.Light
	PrattLamp    *hal.Light
	SaltLamp     *hal.Light

	Onkyo *hal.BinarySensor

	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
	LightsOffTimer hal.Timer
}

type Study struct {
	Lights         *hal.Light       // Main lights
	LuxSensor      *hal.LightSensor // Aqara FP2 (Study)
	PresenceSensor *hal.Entity      // Aqara FP2 (Study)
}

type Upstairs struct {
	LuxSensor      *hal.LightSensor
	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
}

func NewMarnixkade() *Marnixkade {
	cfg, err := hal.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	home := &Marnixkade{
		Connection: hal.NewConnection(*cfg),
		DiningRoom: DiningRoom{
			Lights:         hal.NewLight("light.dining"),
			PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_2"),
		},
		Hallway: Hallway{
			Lights:       hal.NewLight("light.front_hallway"),
			MotionSensor: hal.NewEntity("binary_sensor.hallway_motion"),
		},
		Kitchen: Kitchen{
			PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_4"),
		},
		LivingRoom: LivingRoom{
			ArcherLamp: hal.NewLight("light.archer_lamp"),
			MainLights: hal.LightGroup{
				hal.NewLight("light.hue_filament_bulb_1"),
				hal.NewLight("light.hue_filament_bulb_2"),
				hal.NewLight("light.hue_filament_bulb_3"),
				hal.NewLight("light.hue_filament_bulb_4"),
				hal.NewLight("light.hue_filament_bulb_5"),
			},
			MoroccanLamp: hal.NewLight("light.moroccan_lamp"),
			PrattLamp:    hal.NewLight("light.pratt"),
			SaltLamp:     hal.NewLight("light.salt_lamp"),

			Onkyo: hal.NewBinarySensor("media_player.tx_8270"),

			PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_3"),
		},
		Study: Study{
			Lights:         hal.NewLight("light.study_lights"),
			PresenceSensor: hal.NewEntity("binary_sensor.presence_sensor_fp2_11ad_presence_sensor_1"),
			LuxSensor:      hal.NewLightSensor("sensor.presence_sensor_fp2_11ad_light_sensor_light_level"),
		},
		Upstairs: Upstairs{
			LuxSensor:      hal.NewLightSensor("sensor.presence_sensor_fp2_b6d8_light_sensor_light_level"),
			PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_1"),
		},

		GuestMode: hal.NewBinarySensor("input_boolean.guest_mode"),
	}

	// Walk the struct and find/register all entities
	home.FindEntities(home)

	// Register automations
	home.RegisterAutomations(
		halautomations.NewSensorsTriggerLights().
			WithName("Dining table lights").
			WithSensors(home.DiningRoom.PresenceSensor).
			WithLights(home.DiningRoom.Lights).
			TurnsOffAfter(15*time.Minute),

		halautomations.NewSensorsTriggerLights().
			WithName("Hallway lights").
			WithSensors(home.Hallway.MotionSensor).
			WithLights(home.Hallway.Lights).
			TurnsOffAfter(5*time.Minute),

		hal.NewAutomation().
			WithName("Living room lights").
			WithEntities(home.LivingRoom.PresenceSensor).
			WithAction(func() {
				// Ignore presence changes if someone is actively watching TV or playing music
				if home.LivingRoom.Onkyo.IsOn() {
					return
				}

				if home.LivingRoom.PresenceSensor.IsOn() {
					home.LivingRoom.LightsOffTimer.Cancel()

					// Only turn on the main lights if it's dark
					if home.Upstairs.LuxSensor.Level() < 50 {
						home.LivingRoom.MainLights.TurnOn()
					}

					// Always turn on the archer/pratt lamps
					home.LivingRoom.ArcherLamp.TurnOn()
					home.LivingRoom.PrattLamp.TurnOn()
				} else {
					home.LivingRoom.LightsOffTimer.Start(func() {
						home.LivingRoom.MainLights.TurnOff()
						home.LivingRoom.ArcherLamp.TurnOff()
						home.LivingRoom.PrattLamp.TurnOff()
					}, 15*time.Minute)
				}
			}),

		halautomations.NewSensorsTriggerLights().
			WithName("Study lights").
			WithCondition(func() bool {
				// Disable automation if guest mode is on
				return !home.GuestMode.IsOn()
			}).
			WithSensors(home.Study.PresenceSensor).
			WithLights(home.Study.Lights).
			TurnsOffAfter(5*time.Minute),

		halautomations.NewSensorsTriggerLights().
			WithName("Upstairs bookshelf lamps").
			WithSensors(home.Upstairs.PresenceSensor).
			WithLights(home.LivingRoom.MoroccanLamp, home.LivingRoom.SaltLamp).
			TurnsOffAfter(15*time.Minute),

		halautomations.NewPrintDebug("Debug", hal.NewEntity("event.main_switch_button_1")),
	)

	return home
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	if err := NewMarnixkade().Start(); err != nil {
		slog.Error("Error starting home", "error", err)
		os.Exit(1)
	}

	select {}
}
