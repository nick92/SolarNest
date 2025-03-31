package sensors

import "time"

type Status struct {
	VoltageV      float64   `json:"voltage_v"`
	CurrentA      float64   `json:"current_a"`
	PowerW        float64   `json:"power_w,omitempty"`
	StateOfCharge float64   `json:"soc_percent,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

type Sensor interface {
	GetStatus() (*Status, error)
	GetName() string
}
