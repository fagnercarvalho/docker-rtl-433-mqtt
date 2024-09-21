package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os/exec"

	"github.com/fagnercarvalho/docker-rtl-433-mqtt/mqtt"
)

type Sensor struct {
	Id any `json:"id"`

	// Outdoor temperature/humidity sensor
	Temperature float32 `json:"temperature_C"`
	Humidity    float32 `json:"humidity"`

	// Soil moisture sensor
	Moisture float32 `json:"moisture"`
}

func (s Sensor) IDAsString() string {
	return fmt.Sprint(s.Id)
}

var (
	outdoorSensor = "15909"
	soilSensor    = "0dff63"
)

func main() {
	mqttAddress := flag.String("mqtt-address", "127.0.0.1:1883", "Address + port for MQTT broker to send sensor telemetry")

	flag.Parse()

	client, err := mqtt.NewClient[Sensor](*mqttAddress)
	if err != nil {
		panic(err)
	}

	stdoutPipe := getAntennaStream()

	readFromStream(stdoutPipe, client.PublishMessage)
}

func readFromStream(stdoutPipe io.ReadCloser, publicMqttMessage func(topic string, sensor Sensor) (<-chan error, error)) {
	scanner := bufio.NewScanner(stdoutPipe)

	for scanner.Scan() {
		sensorLine := scanner.Text()

		fmt.Printf("Read sensor from 433 mhz: %v \n", sensorLine)

		var sensor Sensor
		err := json.Unmarshal([]byte(sensorLine), &sensor)
		if err != nil {
			fmt.Printf("Error when trying to unmarshal sensor: %v: %v \n", sensorLine, err)
			continue
		}

		isValid := isValidSensor(sensor.IDAsString())
		if !isValid {
			fmt.Printf("Sensor %v is not valid, skipping \n", sensor.Id)
			continue
		}

		topic := getTopicBySensorID(sensor.IDAsString())

		_, err = publicMqttMessage(topic, sensor)
		if err != nil {
			fmt.Printf("Error when trying to publish message: %v: %v \n", sensor, err)
		}
	}
}

func getAntennaStream() io.ReadCloser {
	// get telemetry from RTL 433 and send to MQTT
	cmd := exec.Command("rtl_433", "-F", "json")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return stdoutPipe
}

func isValidSensor(sensorID string) bool {
	return sensorID == outdoorSensor || sensorID == soilSensor
}

// getTopicBySensorID returns the correct MQTT topic to redirect the sensor data by the sensor ID
// 15909 and 0dff63 are the device IDs for the sensors we are getting the telemetry, change to your IDs
func getTopicBySensorID(id string) string {
	switch id {
	case outdoorSensor:
		return "homeassistant/sensor/balcony/state"
	case soilSensor:
		return "homeassistant/sensor/soil/state"
	}

	return ""
}
