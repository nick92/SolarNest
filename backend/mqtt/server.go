package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttBroker "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/debug"
	"github.com/mochi-mqtt/server/v2/listeners"
	"gorm.io/gorm"
)

type MQTT struct {
	DB *gorm.DB
}

func (mqttMod *MQTT) StartServer() {
	// Create a new MQTT broker instance.
	server := mqttBroker.New(nil)

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	server.AddHook(new(debug.Hook), nil)

	// Create a TCP listener on port 1883.
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "mqtt",
		Address: "0.0.0.0:1883",
	})

	if err := server.AddListener(tcp); err != nil {
		log.Fatalf("Error adding TCP listener: %v", err)
	}

	log.Println("Starting MQTT server on :1883")
	// Serve blocks until error.
	if err := server.Serve(); err != nil {
		log.Fatalf("MQTT server error: %v", err)
	}
}

func (mqttMod *MQTT) StartSubscriber() {
	// Set up the MQTT client options to connect to our broker.
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://0.0.0.0:1883")

	// Create a new MQTT client.
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting MQTT client: %v", token.Error())
	}
	log.Println("MQTT subscriber connected.")

	// Subscribe to the "vedirect/data" topic.
	topic := "vedirect/data"
	if token := client.Subscribe(topic, 0, mqttMod.HandleSubscription); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing: %v", token.Error())
	} else {
		log.Println("Subscription to 'vedirect/data' successful.")
	}

	// Keep the subscriber running.
	select {}
}
