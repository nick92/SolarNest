package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nick92/solarnest/db/schemas"
)

func (mqtt *MQTT) HandleSubscription(client mqtt.Client, msg mqtt.Message) {
	log.Printf("📩 Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))

	var payload schemas.ReadingPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Println("❌ Failed to decode payload:", err)
		return
	}

	var device schemas.Device
	if err := device.RecordMessage(mqtt.DB, &payload, msg.Topic()); err != nil {
		log.Println("❌ DB insert error:", err)
	}
}
