package sensors

import (
	"math/rand"
	"time"
)

type SystemStatus struct {
	SolarPowerW     float64   `json:"solar_power_w"`
	BatteryVoltageV float64   `json:"battery_voltage_v"`
	BatterySOC      float64   `json:"battery_soc"`
	Charging        bool      `json:"charging"`
	Timestamp       time.Time `json:"timestamp"`
}

func GetStatus() SystemStatus {
	rand.Seed(time.Now().UnixNano())

	solarPowerW := rand.Float64()
	batteryVoltageV := rand.Float64()
	batterySOC := rand.Float64()

	return SystemStatus{
		SolarPowerW:     solarPowerW,
		BatteryVoltageV: batteryVoltageV,
		BatterySOC:      batterySOC,
		Charging:        true,
		Timestamp:       time.Now(),
	}
}
