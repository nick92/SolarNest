package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
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

func sendATCommand(port *serial.Port, cmd string, wait time.Duration) (string, error) {
	fullCmd := cmd + "\r\n"
	_, err := port.Write([]byte(fullCmd))
	if err != nil {
		return "", err
	}

	time.Sleep(wait)

	reader := bufio.NewReader(port)
	response := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		response += line
		if line == "\r\n" || line == "OK\r\n" || line == "ERROR\r\n" {
			break
		}
	}
	return response, nil
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Println("Cannot find config.json file, exiting ..")
		return
	}
	c := &serial.Config{Name: cfg.SerialName, Baud: cfg.SerialBaud}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Wake the module
	sendATCommand(port, "AT", 1*time.Second)

	// Optional: set full functionality
	sendATCommand(port, "AT+CFUN=1", 1*time.Second)

	// 1. Set MQTT parameters
	sendATCommand(port, `AT+CMQTTSTART`, 2*time.Second)
	sendATCommand(port, `AT+CMQTTACCQ=0,"client123"`, 1*time.Second)
	sendATCommand(port, `AT+CMQTTCONNECT=0, ${} ,60,1`, 5*time.Second)

	// 2. Set the topic and payload
	topic := "vedirect/data"
	payload := `{"voltage":12.3,"current":5.6}`

	sendATCommand(port, fmt.Sprintf(`AT+CMQTTTOPIC=0,%d`, len(topic)), 1*time.Second)
	sendATCommand(port, topic, 1*time.Second)

	sendATCommand(port, fmt.Sprintf(`AT+CMQTTPAYLOAD=0,%d`, len(payload)), 1*time.Second)
	sendATCommand(port, payload, 1*time.Second)

	// 3. Publish it
	sendATCommand(port, `AT+CMQTTPUB=0,1,60`, 2*time.Second)

	// 4. Close the MQTT session
	sendATCommand(port, `AT+CMQTTDISC=0,60`, 2*time.Second)
	sendATCommand(port, `AT+CMQTTREL=0`, 1*time.Second)
	sendATCommand(port, `AT+CMQTTSTOP`, 1*time.Second)

	log.Println("MQTT message sent!")
}
