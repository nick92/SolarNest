package sensors

import "time"

type SystemStatus struct {
	SolarPowerW     float64   `json:"solar_power_w"`
	BatteryVoltageV float64   `json:"battery_voltage_v"`
	BatterySOC      float64   `json:"battery_soc"`
	Charging        bool      `json:"charging"`
	Timestamp       time.Time `json:"timestamp"`
}

func GetStatus() SystemStatus {
	return SystemStatus{
		SolarPowerW:     123.4,
		BatteryVoltageV: 12.7,
		BatterySOC:      78.9,
		Charging:        true,
		Timestamp:       time.Now(),
	}
}
