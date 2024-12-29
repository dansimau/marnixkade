package main

import (
	"time"

	"github.com/dansimau/hal"
	halautomations "github.com/dansimau/hal/automations"
)

type Study struct {
	Lights       *hal.Light // Main lights
	ClosetLights hal.LightGroup

	ClosetMotionSensor *hal.BinarySensor
	LuxSensor          *hal.LightSensor // Aqara FP2 (Study)
	PresenceSensor     *hal.Entity      // Aqara FP2 (Study)
}

func newStudy() Study {
	return Study{
		Lights: hal.NewLight("light.study_lights"),
		ClosetLights: hal.LightGroup{
			hal.NewLight("light.study_closet_left"),
			hal.NewLight("light.study_closet_right"),
		},

		ClosetMotionSensor: hal.NewBinarySensor("binary_sensor.study_motion_sensor_motion"),
		PresenceSensor:     hal.NewEntity("binary_sensor.presence_sensor_fp2_11ad_presence_sensor_1"),
		LuxSensor:          hal.NewLightSensor("sensor.presence_sensor_fp2_11ad_light_sensor_light_level"),
	}
}

func (s *Study) Automations(home *Marnixkade) []hal.Automation {
	return []hal.Automation{
		halautomations.NewSensorsTriggerLights().
			WithName("Study lights").
			WithCondition(func() bool {
				// Disable automation if guest mode is on
				return !home.GuestMode.IsOn()
			}).
			WithSensors(home.Study.PresenceSensor).
			WithLights(home.Study.Lights).
			TurnsOffAfter(5 * time.Minute),

		halautomations.NewSensorsTriggerLights().
			WithName("Study closet lights").
			WithConditionScene(func() bool { return home.NightMode.IsOn() }, nightLight).
			WithConditionScene(func() bool { return !home.NightMode.IsOn() }, brightLight).
			WithSensors(home.Study.ClosetMotionSensor).
			WithLights(home.Study.ClosetLights).
			TurnsOffAfter(1 * time.Minute),
	}
}
