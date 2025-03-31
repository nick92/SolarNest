package sensors

import (
	"time"

	"github.com/d2r2/go-i2c"
)

type INA219Sensor struct{}

func (s *INA219Sensor) GetName() string {
	return "INA219"
}

func (s *INA219Sensor) GetStatus() (*Status, error) {
	i2c, err := i2c.NewI2C(0x40, 1)

	if err != nil {
		return nil, err
	}
	defer i2c.Close()

	rawBus, _, err := i2c.ReadRegBytes(0x02, 2)

	if err != nil {
		return nil, err
	}

	busRaw := (uint16(rawBus[0]) << 8) | uint16(rawBus[1])
	busVoltage := float64((busRaw>>3)*4) / 1000.0 // mV to V

	// Read Current Register (0x04, 2)
	rawCurrent, _, err := i2c.ReadRegBytes(0x04, 2)

	if err != nil {
		return nil, err
	}

	currentRaw := int16(rawCurrent[0])<<8 | int16(rawCurrent[1])
	current := float64(currentRaw) / 1000.0 // mA to A

	return &Status{
		VoltageV:  busVoltage,
		CurrentA:  current,
		PowerW:    busVoltage * current,
		Timestamp: time.Now(),
	}, nil
}
