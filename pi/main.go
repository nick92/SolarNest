package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tarm/serial"
)

func loadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home dir: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "solarnest", "config.json")

	fmt.Println("📄 Reading config from:", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func sendATCommand(port *serial.Port, cmd string, wait time.Duration) (string, error) {
	log.Println("Sending AT Command: " + cmd)

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
			log.Println("Error on AT Commend:" + err.Error())
			break
		}

		line = strings.TrimSpace(line)
		if line != "" {
			log.Println("⬅️  ", line)
			response += line + "\n"

			// Exit early if we see a final response
			if strings.Contains(line, "OK") || strings.Contains(line, "ERROR") {
				break
			}
		}
	}
	return response, nil
}
func sendATCommandWithPayload(port *serial.Port, cmd string, payload string, wait time.Duration) (string, error) {
	log.Println("➡️ Sending AT Command: " + cmd)

	fullCmd := cmd + "\r\n"
	_, err := port.Write([]byte(fullCmd))
	if err != nil {
		log.Println("❌ Error writing command:", err)
		return "", err
	}

	time.Sleep(wait)

	// Read until we get a '>' (or ERROR)
	reader := bufio.NewReader(port)
	var response string

	for {
		b, err := reader.ReadByte()
		if err != nil {
			log.Println("❌ Read error:", err)
			break
		}
		ch := string(b)
		response += ch

		if ch == ">" {
			log.Println("✅ Prompt received: >", "→ writing payload:", payload)
			time.Sleep(200 * time.Millisecond)

			// Write the payload with NO newline
			_, err = port.Write([]byte(payload))
			if err != nil {
				log.Println("❌ Error writing payload:", err)
				return response, err
			}

			time.Sleep(wait)
			break
		}

		if strings.Contains(response, "ERROR") {
			log.Println("❌ ERROR received before payload prompt")
			return response, fmt.Errorf("modem returned ERROR")
		}

		if strings.Contains(response, "OK") {
			break
		}
	}

	// Read final response after payload
	reader = bufio.NewReader(port)

	payloadResponse := ""

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Println("❌ Error reading final response:", err)
		return response, err
	}

	payloadResponse += line

	log.Println("Payload Response: " + payloadResponse)

	final := strings.TrimSpace(payloadResponse)
	log.Println("⬅️  Final response:", final)
	response += final

	return response, nil
}

func generateRandomPayload() string {
	rand.Seed(time.Now().UnixNano())

	voltage := 11.0 + rand.Float64()*2.0 // 11.0V – 13.0V
	current := 0.0 + rand.Float64()*10.0 // 0A – 10A

	payload := fmt.Sprintf(`{"voltage":%.2f,"current":%.2f}`, voltage, current)
	return payload
}

func sendReading(port *serial.Port, cfg *Config) {
	// Wake the SIM7028
	sendATCommand(port, "AT", 2*time.Second)
	sendATCommand(port, "AT+CFUN=1", 1*time.Second)

	// 1. Start MQTT session and set client ID
	sendATCommand(port, `AT+CMQTTSTART`, 2*time.Second)
	sendATCommand(port, `AT+CMQTTACCQ=0,"mqtt"`, 1*time.Second)

	sendATCommand(port, fmt.Sprintf(`AT+CMQTTCONNECT=0,"%s",60,1`, cfg.MQTTServerURL), 5*time.Second)

	// 3. Set topic and payload
	topic := cfg.MQTTTopic
	payload := generateRandomPayload()

	sendATCommandWithPayload(port, fmt.Sprintf(`AT+CMQTTTOPIC=0,%d`, len(topic)), topic, 2*time.Second)

	sendATCommandWithPayload(port, fmt.Sprintf(`AT+CMQTTPAYLOAD=0,%d`, len(payload)), payload, 2*time.Second)

	// 4. Publish the message
	sendATCommand(port, `AT+CMQTTPUB=0,1,60`, 2*time.Second)

	// 5. Gracefully disconnect
	sendATCommand(port, `AT+CMQTTDISC=0,60`, 2*time.Second)
	sendATCommand(port, `AT+CMQTTREL=0`, 1*time.Second)
	sendATCommand(port, `AT+CMQTTSTOP`, 1*time.Second)

	log.Println("✅ MQTT message sent from SIM7028!")
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Println("Cannot find config.json file, exiting ..")
		return
	}

	log.Println("Connecting to serial port SIM7028!")

	c := &serial.Config{Name: cfg.SerialName, Baud: cfg.SerialBaud}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	log.Println("✅ Serial port connected!")

	for {
		time.Sleep(time.Duration(cfg.IntervalMinutes) * time.Second)
		sendReading(port, cfg) // your logic that connects, publishes, disconnects
	}
}
