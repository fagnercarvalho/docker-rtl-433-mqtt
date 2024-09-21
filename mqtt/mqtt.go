package mqtt

import (
	"errors"
	"fmt"
	"net"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewClient(mqttAddress string, mqttPort int) (mqtt.Client, error) {
	if mqttAddress == "" || mqttPort == 0 {
		return nil, errors.New("MQTT address and port must be provided")
	}

	client := newClient(mqttAddress, mqttPort)

	err := connectToBroker(client)
	if err != nil {
		return nil, fmt.Errorf("error to connect to MQTT broker: %w", err)
	}

	return client, nil
}

func PublishMessage[T any](client mqtt.Client) func(topic string, message T) {
	return func(topic string, message T) {
		t := client.Publish(topic, 0, false, message)
		go func() {
			<-t.Done()
			if t.Error() != nil {
				fmt.Printf("Error when trying to publish message %v to topic %v: %v", message, topic, t.Error())
			}
		}()
	}
}

func newClient(address string, port int) mqtt.Client {
	clientOptions := mqtt.NewClientOptions()
	clientOptions.Dialer = &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(address),
			Port: port,
		},
	}

	return mqtt.NewClient(clientOptions)
}

func connectToBroker(client mqtt.Client) error {
	t := client.Connect()

	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}
