package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Bedroom struct {
	AllLights *hal.Light

	MainLights    *hal.Light
	GoldenSunLamp *hal.Light
	ClosetLights  hal.LightGroup
	BedLights     *hal.Light

	ClosetMotionSensor *hal.BinarySensor // Hue Motion Sensor
	PresenceSensor     *hal.BinarySensor // Aqara FP2 (Bedroom)
}

func newBedroom() Bedroom {
	return Bedroom{
		AllLights:     hal.NewLight("light.bedroom"),
		MainLights:    hal.NewLight("light.bedroom_lights"),
		GoldenSunLamp: hal.NewLight("light.golden_sun"),
		ClosetLights: hal.LightGroup{
			hal.NewLight("light.bedroom_closet_left"),
			hal.NewLight("light.bedroom_closet_right"),
		},
		BedLights: hal.NewLight("light.bed_strip"),

		ClosetMotionSensor: hal.NewBinarySensor("binary_sensor.bedroom_motion"),
		PresenceSensor:     hal.NewBinarySensor("binary_sensor.presence_sensor_fp2_1a4f_presence_sensor_1"),
	}
}

func (b *Bedroom) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Bedroom lights").
			WithCondition(func() bool {
				return !home.NightMode.IsOn() // Don't auto turn on lights if night mode is on
			}).
			WithSensors(home.Bedroom.PresenceSensor).
			TurnsOnLights(
				home.Bedroom.MainLights,
				home.Bedroom.GoldenSunLamp,
				home.Bedroom.BedLights,
			).
			TurnsOffLights(home.Bedroom.AllLights).
			TurnsOffAfter(5 * time.Minute),

		halautomations.NewSensorsTriggerLights().
			WithName("Bedroom closet lights").
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, nightLight).
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			WithSensors(home.Bedroom.ClosetMotionSensor).
			WithLights(home.Bedroom.ClosetLights).
			TurnsOffAfter(30 * time.Second),

		halautomations.NewTimer("Detect person in bed").
			WithEntities(home.Bedroom.PresenceSensor).
			Condition(home.Bedroom.PresenceSensor.IsOn).
			Duration(15 * time.Minute).
			Run(func() {
				if home.Bedroom.PresenceSensor.IsOn() {
					home.NightMode.TurnOn()
				}
			}),

		// halautomations.NewTimer("Detect everyone out of bed").
		// 	WithEntities(home.Bedroom.PresenceSensor).
		// 	Condition(home.Bedroom.PresenceSensor.IsOff).
		// 	Duration(15 * time.Minute).
		// 	Run(func() {
		// 		if home.Bedroom.PresenceSensor.IsOff() {
		// 			home.NightMode.TurnOff()
		// 		}
		// 	}),
	}
}
