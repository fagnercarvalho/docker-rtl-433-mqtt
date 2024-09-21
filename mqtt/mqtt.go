package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client[T any] interface {
	PublishMessage(topic string, message T) (<-chan error, error)
	Subscribe(topic string) (<-chan string, <-chan error)
}

type client[T any] struct {
	internalClient mqtt.Client
}

func NewClient[T any](host string, port string) (Client[T], error) {
	if host == "" || port == "" {
		return nil, errors.New("MQTT host and port must be provided")
	}

	internalClient, err := newClient(host, port)
	if err != nil {
		return nil, fmt.Errorf("error to instantiate new MQTT client: %w", err)
	}

	err = connectToBroker(internalClient)
	if err != nil {
		return nil, fmt.Errorf("error to connect to MQTT broker: %w", err)
	}

	return client[T]{internalClient: internalClient}, nil
}

func (c client[T]) PublishMessage(topic string, message T) (<-chan error, error) {
	channel := make(chan error)

	bytes, err := json.Marshal(message)
	if err != nil {
		return channel, err
	}

	t := c.internalClient.Publish(topic, 0, false, bytes)
	go func() {
		<-t.Done()
		if t.Error() != nil {
			fmt.Printf("Error when trying to publish message %v to topic %v: %v", message, topic, t.Error())
			channel <- t.Error()
		}

		channel <- nil
	}()

	return channel, nil
}

func (c client[T]) Subscribe(topic string) (<-chan string, <-chan error) {
	errCh := make(chan error)
	successCh := make(chan string)

	t := c.internalClient.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
		successCh <- string(msg.Payload())
	})
	go func() {
		<-t.Done()
		if t.Error() != nil {
			fmt.Printf("Error when trying to subscribe to topic %v: %v", topic, t.Error())
			errCh <- t.Error()
		}

		errCh <- nil
	}()

	return successCh, errCh
}

func newClient(host, port string) (mqtt.Client, error) {
	clientOptions := mqtt.NewClientOptions()

	parsedURL, err := url.Parse("mqtt://" + host + ":" + port)
	if err != nil {
		return nil, err
	}

	clientOptions.Servers = []*url.URL{parsedURL}

	return mqtt.NewClient(clientOptions), nil
}

func connectToBroker(client mqtt.Client) error {
	t := client.Connect()

	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}
