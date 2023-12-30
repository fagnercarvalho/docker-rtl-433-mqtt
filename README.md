## docker-rtl-433-mqtt

## TL;DR

This is an app that runs [RTL_433](https://github.com/merbanan/rtl_433) to get data from Radio signals and flow to MQTT via the [Mosquitto](https://mosquitto.org/) broker.

## What this does

This is a Docker container that:
- Send MQTT messages to create sensors in Home Assistant
- Start reading Radio signals from an [RTL-SDR](https://en.wikipedia.org/wiki/Software-defined_radio) (Realtek Software Defined Radio) antenna using [RTL_433](https://github.com/merbanan/rtl_433)
- Send messages to MQTT via `mosquitto_pub` MQTT client in a format that can be read by Home Assistant

## Prerequisites

You will need:
- Docker
- An RTL-SDR USB dongle. I use the `NooElec NESDR Mini USB RTL-SDR` but [RTL_433](https://github.com/merbanan/rtl_433) supports a lot of different models, choose the one that you prefer
- Home Assistant or something on the other side to consume the MQTT messages

## Running

To run this:

1. Create an `.env` file like this.
   This will be used by container to connect to MQTT.
```
MQTT_HOST=<value>
MQTT_PORT=<value>
```

2. Expose the correct RTL-SDR USB dongle to the container
```yaml
devices:
  - /dev/bus/usb/001/005:/dev/bus/usb/001/005
```

3. Run Docker Compose
```shell
docker compose up --build
```

If everything goes well you will see the sensor states in Home Assistant like this:

![Home Assistant](home-assistant.png)