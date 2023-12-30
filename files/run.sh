#!/bin/bash

# create MQTT topics for sensors
# change to your sensors
mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/outdoor/temperature/config" -f mqtt-temperature-config-message
mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/outdoor/humidity/config" -f mqtt-humidity-config-message

# get telemetry from RTL 433 and send to MQTT
# 15909 is the device ID for the sensor we are getting the telemetry, change to your ID
rtl_433 -F json | jq -c 'select(.id == 15909)' | mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/balcony/state" -l

# uncomment to test
#cat mock-telemetry | jq -c 'select(.id == 15909)' | mosquitto_pub -h $MQTT_HOST -p $MQTT_PORT -t "homeassistant/sensor/balcony/state" -l
