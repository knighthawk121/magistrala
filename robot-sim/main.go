package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Velocity    float64 `json:"velocity"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	Timestamp   int64   `json:"timestamp"`
}

func main() {
	credsPath := "/home/khawk/rodrop/magistrala/NGS-Default-magistrala-service.creds"
	opts := []nats.Option{
		nats.UserCredentials(credsPath),
	}

	nc, err := nats.Connect("tls://connect.ngs.global:4222", opts...)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to Synadia NATS successfully.")

	// A common pattern for magistrala channels: channels.<channel_id>.messages.*
	topic := "channels.robot.messages.senml"
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Publishing simulated sensor data to topic '%s' every 5s...", topic)

	for {
		select {
		case <-ticker.C:
			data := SensorData{
				Temperature: 20.0 + rand.Float64()*10.0,
				Velocity:    0.5 + rand.Float64()*1.5,
				PosX:        rand.Float64() * 100,
				PosY:        rand.Float64() * 100,
				Timestamp:   time.Now().UnixNano(),
			}
			payload, _ := json.Marshal(data)

			err := nc.Publish(topic, payload)
			if err != nil {
				log.Printf("Failed to publish: %v", err)
			} else {
				log.Printf("Published: %s", string(payload))
			}
		case <-sigChan:
			log.Println("Shutting down simulator...")
			return
		}
	}
}
