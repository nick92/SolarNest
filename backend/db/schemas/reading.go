package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Reading struct {
	ID         uint `gorm:"primaryKey"`
	Voltage    float64
	Current    float64
	Topic      string
	DeviceID   *uint
	ReceivedAt time.Time `gorm:"autoCreateTime"`
}

func (reading *Reading) Insert(db *gorm.DB) error {
	if err := db.Create(&reading).Error; err != nil {
		return err
	}

	return nil
}
