package mqtt

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Test struct {
	Testing     string `json:"testing"`
	OnlyTesting int    `json:"only_testing"`
}

func TestPublishMessage(t *testing.T) {
	// instantiate MQTT client
	client, err := NewClient[Test]("127.0.0.1", "1883")
	require.NoError(t, err)

	// subscribe to topic
	successChannel, errChannel := client.Subscribe("test-topic")

	select {
	case err := <-errChannel:
		require.NoError(t, err)
	case <-time.After(time.Second * 2):
		assert.Fail(t, "Waited for too long for MQTT subscription")
	}

	payload := Test{Testing: "1", OnlyTesting: 3}

	// assert message
	var gotMessage bool
	go func() {
		select {
		case message := <-successChannel:
			var actualPayload Test
			err := json.Unmarshal([]byte(message), &actualPayload)
			require.NoError(t, err)

			assert.Equal(t, payload, actualPayload)
		case <-time.After(time.Second * 2):
			assert.Fail(t, "Waited for too long for MQTT message")
		}

		gotMessage = true
	}()

	// publish message
	errChannel, err = client.PublishMessage("test-topic", payload)
	require.NoError(t, err)

	select {
	case err := <-errChannel:
		require.NoError(t, err)
	case <-time.After(time.Second * 2):
		assert.Fail(t, "Waited for too long for MQTT successful response")
	}

	// assert that got message
	assert.Eventually(t, func() bool {
		return gotMessage
	}, time.Second*5, time.Millisecond*500)
}
