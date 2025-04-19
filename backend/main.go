package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nick92/solarnest/api"
	"github.com/nick92/solarnest/mqtt"
)

func main() {
	go initGinServer()
	initMQTTServer()
}

func initGinServer() {
	r := gin.Default()

	api.SetupRoutes(r)

	log.Println("HTTP server running on :8080")
	r.Run("0.0.0.0:8080")

}

func initMQTTServer() {

	// Start the MQTT broker in a separate goroutine.
	go mqtt.StartServer()

	// Give the broker a couple of seconds to fully start.
	time.Sleep(2 * time.Second)

	// Start the MQTT subscriber (client) that will read messages from our broker.
	mqtt.StartSubscriber()

}
