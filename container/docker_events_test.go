package container

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/logocomune/webhookdocker/message"
	"github.com/stretchr/testify/assert"
)

func TestParseEvent(t *testing.T) {
	// Test container event
	t.Run("Container event", func(t *testing.T) {
		// Create a sample Docker container event
		dockerEvent := events.Message{
			Type:   "container",
			Action: "start",
			Actor: events.Actor{
				ID: "container123",
				Attributes: map[string]string{
					"name":     "test-container",
					"image":    "test-image",
					"signal":   "SIGTERM",
					"exitCode": "0",
				},
			},
			Time: 1609459200, // 2021-01-01 00:00:00 UTC
		}

		// Parse the event
		event := parseEvent(dockerEvent)

		// Verify the parsed event
		assert.Equal(t, "container", event.MetaData.Type)
		assert.Equal(t, "start", event.MetaData.Action)
		assert.Equal(t, time.Unix(1609459200, 0), event.MetaData.Time)
		assert.Equal(t, "container123", event.Container.ID)
		assert.Equal(t, "test-container", event.Container.Name)
		assert.Equal(t, "test-image", event.Container.Image)
		assert.Equal(t, "SIGTERM", event.ContainerStatus.Signal)
		assert.Equal(t, "0", event.ContainerStatus.ExitCode)
	})

	// Test volume event
	t.Run("Volume event", func(t *testing.T) {
		// Create a sample Docker volume event
		dockerEvent := events.Message{
			Type:   "volume",
			Action: "create",
			Actor: events.Actor{
				ID: "volume123",
				Attributes: map[string]string{
					"destination": "/data",
					"container":   "container456",
				},
			},
			Time: 1609459200, // 2021-01-01 00:00:00 UTC
		}

		// Parse the event
		event := parseEvent(dockerEvent)

		// Verify the parsed event
		assert.Equal(t, "volume", event.MetaData.Type)
		assert.Equal(t, "create", event.MetaData.Action)
		assert.Equal(t, time.Unix(1609459200, 0), event.MetaData.Time)
		assert.Equal(t, "volume123", event.Volume.ID)
		assert.Equal(t, "/data", event.Volume.Destination)
		assert.Equal(t, "container456", event.Container.ID)
	})

	// Test network event
	t.Run("Network event", func(t *testing.T) {
		// Create a sample Docker network event
		dockerEvent := events.Message{
			Type:   "network",
			Action: "connect",
			Actor: events.Actor{
				ID: "network123",
				Attributes: map[string]string{
					"name":      "test-network",
					"type":      "bridge",
					"container": "container789",
				},
			},
			Time: 1609459200, // 2021-01-01 00:00:00 UTC
		}

		// Parse the event
		event := parseEvent(dockerEvent)

		// Verify the parsed event
		assert.Equal(t, "network", event.MetaData.Type)
		assert.Equal(t, "connect", event.MetaData.Action)
		assert.Equal(t, time.Unix(1609459200, 0), event.MetaData.Time)
		assert.Equal(t, "network123", event.Network.ID)
		assert.Equal(t, "test-network", event.Network.Name)
		assert.Equal(t, "bridge", event.Network.Type)
		assert.Equal(t, "container789", event.Container.ID)
	})
}

func TestContextDoneError(t *testing.T) {
	err := &ContextDoneError{}
	assert.Equal(t, "DockerEvents: Context done. Exit", err.Error())
}

func TestDockerEventsContextCancellation(t *testing.T) {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to receive events
	eventChan := make(chan message.Event)

	// Create a configuration
	cfg := DockerCfg{
		ContainerEvents:  true,
		VolumeEvents:     true,
		NetworkEvents:    true,
		ContainerActions: []string{"start", "stop"},
		NetworkActions:   []string{"create", "destroy"},
		VolumeActions:    []string{"create", "destroy"},
	}

	// Cancel the context immediately to simulate context cancellation
	cancel()

	// Call DockerEvents with the cancelled context
	err := DockerEvents(ctx, eventChan, cfg, "test-version", "test-date")

	// Verify that we get a ContextDoneError
	assert.IsType(t, &ContextDoneError{}, err)
	assert.Equal(t, "DockerEvents: Context done. Exit", err.Error())
}
