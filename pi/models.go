package main

type Config struct {
	ServerURL       string `json:"server_url"`
	ClientID        string `json:"client_id"`
	SerialName      string `json:"serial_name"`
	SerialBaud      int    `json:"serial_baud"`
	IntervalMinutes int    `json:"interval_minutes"`
	AuthToken       string `json:"auth_token"`
}

type ServerPingData struct {
	Message string `json:"message"`
}
