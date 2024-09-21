package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFromStream(t *testing.T) {
	b, err := os.ReadFile("files/mock-telemetry")
	if err != nil {
		t.Fatalf("Error while reading file: %v", err)
	}
	require.NoError(t, err)

	stream := io.NopCloser(bytes.NewReader(b))

	expectedSensors := []Sensor{
		{Id: float64(15909), Temperature: 66.700, Humidity: 50, Moisture: 0},
		{Id: float64(15909), Temperature: 66.700, Humidity: 50, Moisture: 0},
		{Id: float64(15909), Temperature: 66.700, Humidity: 50, Moisture: 0},
		{Id: "0dff63", Temperature: 0, Humidity: 0, Moisture: 52},
	}

	var sensors []Sensor
	readFromStream(stream, func(_ string, sensor Sensor) {
		sensors = append(sensors, sensor)
	})

	assert.ElementsMatch(t, sensors, expectedSensors)
}
