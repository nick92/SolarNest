package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func main() {
	cfg, err := loadConfig("config.json")
	ticker := time.NewTicker(time.Duration(cfg.IntervalMinutes) * time.Minute)
	defer ticker.Stop()

	if err != nil {
		log.Println("Cannot find config.json file, exiting ..")
		return
	}

	log.Println("Starting Pi heartbeat service...")

	pingResponse := sendPingMessage(cfg)

	if pingResponse != nil && pingResponse.Message == "pong" {
		for t := range ticker.C {
			log.Println("Ticker at: ", t)
			sendMessage(cfg)
		}
	}
}
