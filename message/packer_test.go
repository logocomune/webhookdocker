package message

import (
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
)

func TestPacker_EventPacker(t *testing.T) {
	l, _ := lru.New(1)
	packer := Packer{
		cache:    l,
		nodeName: "Test",
	}

	events := make([]Event, 2, 2)
	events[0] = Event{
		MetaData: MetaData{
			Type:   dockerWebhook,
			Action: "",
			Time:   time.Now(),
		},
		StartupInfo:     StartupInfo{},
		Container:       Container{},
		ContainerStatus: ContainerStatus{},
		Volume:          Volume{},
		Network:         Network{},
	}

	events[1] = Event{
		MetaData: MetaData{
			Type:   "container",
			Action: "create",
			Time:   time.Now(),
		},
		StartupInfo: StartupInfo{},
		Container: Container{
			ID:    "abcd",
			Name:  "ContainerName",
			Image: "Image",
		},
		ContainerStatus: ContainerStatus{},
		Volume:          Volume{},
		Network:         Network{},
	}
	eventPacker := packer.EventPacker(events)
	assert.NotEmpty(t, eventPacker)
	assert.Equal(t, 2, len(eventPacker))
	assert.NotEmpty(t, eventPacker[dockerWebhook])

	assert.NotEmpty(t, eventPacker["abcd"])
	assert.Equal(t, Container{
		ID:    "abcd",
		Name:  "ContainerName",
		Image: "Image",
	},
		eventPacker["abcd"].Container)

	assert.Equal(t,1,len(eventPacker["abcd"].Events))

	assert.Equal(t,eventPacker["abcd"].Events[0],events[1])

}
