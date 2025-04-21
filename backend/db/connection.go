package db

import (
	"log"
	"os"

	"github.com/nick92/solarnest/db/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&schemas.Device{}, &schemas.Reading{}); err != nil {
		return nil, err
	}

	log.Println("✅ GORM DB ready")
	return db, nil
}
