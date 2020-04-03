package container

import (
	"regexp"
	"strings"

	"docker.io/go-docker/api/types/events"
)

type filter struct {
	actions       map[string]struct{}
	imageName     *regexp.Regexp
	containerName *regexp.Regexp
}

const (
	containerEvent = events.ContainerEventType + "_"
	networkEvent   = events.NetworkEventType + "_"
	volumeEvent    = events.VolumeEventType + "_"
)

func newFilter(cfg DockerCfg) (*filter, error) {
	actions := make(map[string]struct{})

	var imageName, containerName *regexp.Regexp

	var err error

	if cfg.FilterImage != "" {
		imageName, err = regexp.Compile(cfg.FilterImage)
		if err != nil {
			return nil, err
		}
	}

	if cfg.FilterName != "" {
		containerName, err = regexp.Compile(cfg.FilterName)
		if err != nil {
			return nil, err
		}
	}

	for _, s := range cfg.ContainerActions {
		actions[containerEvent+s] = struct{}{}
	}

	for _, s := range cfg.NetworkActions {
		actions[networkEvent+s] = struct{}{}
	}

	for _, s := range cfg.VolumeActions {
		actions[volumeEvent+s] = struct{}{}
	}

	return &filter{
		actions:       actions,
		imageName:     imageName,
		containerName: containerName,
	}, nil
}

func (f *filter) accept(event events.Message) bool {

	typeEvent := ""

	switch event.Type {
	case events.ContainerEventType:
		if f.containerName != nil {
			containerName := event.Actor.Attributes["name"]
			if !f.containerName.MatchString(containerName) {
				return false
			}
		}

		if f.imageName != nil {
			imageName := event.Actor.Attributes["image"]
			if !f.imageName.MatchString(imageName) {
				return false
			}
		}

		typeEvent = containerEvent
	case events.NetworkEventType:
		typeEvent = networkEvent
	case events.VolumeEventType:
		typeEvent = volumeEvent
	}

	action := strings.Split(event.Action, ":")

	typeEvent += action[0]

	_, ok := f.actions[typeEvent]

	return ok
}
