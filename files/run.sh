#!/bin/bash

echo "Running run.sh"

# create MQTT topics for sensors
# change to your sensors
mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/outdoor/temperature/config" -f mqtt-temperature-config-message
mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/outdoor/humidity/config" -f mqtt-humidity-config-message
mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/soil/moisture/config" -f mqtt-moisture-config-message
