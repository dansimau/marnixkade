package main

import (
	"context"
	"time"

	"github.com/dansimau/hal"
)

type LivingRoom struct {
	ArcherLamp   *hal.Light
	MainLights   hal.LightGroup // Hue Filament Bulbs
	MoroccanLamp *hal.Light
	PrattLamp    *hal.Light
	SaltLamp     *hal.Light

	TV *hal.BinarySensor

	PresenceSensor *hal.BinarySensor // Aqara FP2 (Bar)
	LightsOffTimer hal.Timer
}

func newLivingRoom() LivingRoom {
	return LivingRoom{
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

		TV: hal.NewBinarySensor("media_player.lounge"),

		PresenceSensor: hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_b6d8_presence_sensor_3"),
	}
}

func (l *LivingRoom) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		hal.NewAutomation().
			WithName("Living room lights").
			WithEntities(home.LivingRoom.PresenceSensor).
			WithAction(func(ctx context.Context, _ hal.EntityInterface) {
				// Ignore presence changes if someone is actively watching TV or playing music
				if home.LivingRoom.TV.GetState().State == "playing" {
					return
				}

				if home.LivingRoom.PresenceSensor.IsOn() {
					home.LivingRoom.LightsOffTimer.Cancel()

					// Only turn on the main lights if it's dark
					if home.Upstairs.LuxSensor.Level() < 50 {
						home.LivingRoom.MainLights.TurnOnContext(ctx)
					}

					// Always turn on the archer/pratt lamps
					home.LivingRoom.ArcherLamp.TurnOnContext(ctx)
					home.LivingRoom.PrattLamp.TurnOnContext(ctx)
				} else {
					home.LivingRoom.LightsOffTimer.StartContext(ctx, func(ctx context.Context) {
						home.LivingRoom.MainLights.TurnOffContext(ctx)
						home.LivingRoom.ArcherLamp.TurnOffContext(ctx)
						home.LivingRoom.PrattLamp.TurnOffContext(ctx)
					}, 15*time.Minute)
				}
			}),
		hal.NewAutomation().
			WithName("Turn off main lights when TV starts").
			WithEntities(home.LivingRoom.TV).
			WithAction(func(ctx context.Context, _ hal.EntityInterface) {
				if home.LivingRoom.TV.GetState().State != "playing" {
					return
				}
				home.LivingRoom.MainLights.TurnOffContext(ctx)
			}),
	}
}
