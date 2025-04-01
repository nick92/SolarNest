package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func sendPingMessage(cfg *Config) *ServerPingData {

	response, err := http.Get(cfg.ServerURL + "/api/ping")
	if err != nil {
		log.Println("Error sending POST request:", err)
		return nil
	}
	defer response.Body.Close()

	responseData, err := decodeResponse[ServerPingData](response.Body)

	if err != nil {
		log.Println("Error decoding response data:", err)
		return nil
	}

	return &responseData.Data
}

func decodeResponse[T any](body io.Reader) (*ServerResponse[T], error) {
	var response ServerResponse[T]
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func sendMessage(cfg *Config) {
	payload := PostPayload{
		Message: "Hello from Raspberry Pi!",
		Time:    time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	response, err := http.Post(cfg.ServerURL+"/api/ping", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer response.Body.Close()

	log.Println("Sent message, server responded with status:", response.Status)
}
