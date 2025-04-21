package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nick92/solarnest/api"
	"github.com/nick92/solarnest/db"
	"github.com/nick92/solarnest/mqtt"
	"gorm.io/gorm"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	go initGinServer(db)
	initMQTTServer(db)
}

func initGinServer(db *gorm.DB) {
	r := gin.Default()

	api.SetupRoutes(r, db)

	log.Println("HTTP server running on :8080")
	r.Run("0.0.0.0:8080")
}

func initMQTTServer(db *gorm.DB) {
	mqtt := &mqtt.MQTT{DB: db}

	// Start the MQTT broker in a separate goroutine.
	go mqtt.StartServer()

	// Give the broker a couple of seconds to fully start.
	time.Sleep(2 * time.Second)

	// Start the MQTT subscriber (client) that will read messages from our broker.
	mqtt.StartSubscriber()
}
