package container

import (
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/events"
)

type filter struct {
	actions             map[string]struct{}
	imageName           *regexp.Regexp
	containerName       *regexp.Regexp
	invertImageName     bool
	invertContainerName bool
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
		actions[string(containerEvent)+s] = struct{}{}
	}

	for _, s := range cfg.NetworkActions {
		actions[string(networkEvent)+s] = struct{}{}
	}

	for _, s := range cfg.VolumeActions {
		actions[string(volumeEvent)+s] = struct{}{}
	}

	return &filter{
		actions:             actions,
		imageName:           imageName,
		containerName:       containerName,
		invertContainerName: cfg.NegateFilterName,
		invertImageName:     cfg.NegateFilterImage,
	}, nil
}

func (f *filter) accept(event events.Message) bool {
	typeEvent := ""

	switch event.Type {
	case events.ContainerEventType:
		if f.containerName != nil {
			cN := event.Actor.Attributes["name"]
			match := f.containerName.MatchString(cN)

			if match == f.invertContainerName {
				return false
			}
		}

		if f.imageName != nil {
			iN := event.Actor.Attributes["image"]
			match := f.imageName.MatchString(iN)

			if match == f.invertImageName {
				return false
			}
		}

		typeEvent = string(containerEvent)
	case events.NetworkEventType:
		typeEvent = string(networkEvent)
	case events.VolumeEventType:
		typeEvent = string(volumeEvent)
	}

	action := strings.Split(string(event.Action), ":")

	typeEvent += action[0]

	_, ok := f.actions[typeEvent]

	return ok
}
