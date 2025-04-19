package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttBroker "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/debug"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func StartServer() {
	// Create a new MQTT broker instance.
	server := mqttBroker.New(nil)

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	server.AddHook(new(debug.Hook), nil)

	// config := listeners.Config{
	// 	Type:    "mqtt",
	// 	Address: "0.0.0.0:1883",
	// }

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

func StartSubscriber() {
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
	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing: %v", token.Error())
	} else {
		log.Println("Subscription to 'vedirect/data' successful.")
	}

	// Keep the subscriber running.
	select {}
}
