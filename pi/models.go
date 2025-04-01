package main

import "time"

type PostPayload struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

type Config struct {
	ServerURL       string `json:"server_url"`
	IntervalMinutes int    `json:"interval_minutes"`
	AuthToken       string `json:"auth_token"`
}

type ServerResponse[T any] struct {
	Status       string    `json:"status"`
	Data         T         `json:"data"`
	ErrorCode    int       `json:"code"`
	ErrorMessage string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
}

type ServerPingData struct {
	Message string `json:"message"`
}
