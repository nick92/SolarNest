package schemas

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID        uint   `gorm:"primaryKey"`
	UID       string `gorm:"uniqueIndex"` // e.g., "pi-zero-001" or SIM ICCID
	Name      string `gorm:"size:255"`
	Status    string `gorm:"size:100"`
	Location  string `gorm:"size:255"`
	LastSeen  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Readings []Reading `gorm:"foreignKey:DeviceID"` // optional relation
}

func (device *Device) GetDevices(db *gorm.DB) []Device {
	var devices []Device

	if err := db.Order("last_seen desc").Find(&devices).Error; err != nil {
		return devices
	}

	return devices
}

func (device *Device) RecordMessage(db *gorm.DB, payload *ReadingPayload, topic string) error {
	// Look for device by UID
	result := db.Where("uid = ?", payload.DeviceID).First(&device)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Not found — create it
			device.UID = payload.DeviceID
			device.Status = "online"
			if err := db.Create(&device).Error; err != nil {
				return err
			}
			log.Printf("🆕 Registered new device: %s", device.UID)
		} else {
			return result.Error
		}
	}

	// Insert new reading
	reading := Reading{
		Voltage:  payload.Voltage,
		Current:  payload.Current,
		Topic:    topic,
		DeviceID: &device.ID,
	}

	reading.Insert(db)

	log.Printf("✅ Saved reading for device %s (ID %d)", device.UID, device.ID)
	return nil
}
