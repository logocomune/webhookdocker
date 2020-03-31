package container

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/events"
	"docker.io/go-docker/api/types/filters"
	"github.com/logocomune/webhookdocker/message"
)

type DockerCfg struct {
	ContainerEvents bool
	VolumeEvents    bool
	NetworkEvents   bool
	ShowRunning     bool
}

type ContextDoneError struct {
}

func (c *ContextDoneError) Error() string { return "DockerEvents: Context done. Exit" }

func DockerEvents(ctx context.Context, cEvnt chan message.Event, cfg DockerCfg) error {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})

	if err != nil {
		return err
	}

	var v types.Version

	if v, err = cli.ServerVersion(ctx); err != nil {
		return err
	}
	log.Println("Docker version", v.Version, "Api version", v.APIVersion)

	cEvnt <- message.Event{
		MetaData: message.MetaData{
			Type:   "docker-webhook",
			Action: "startup",
			Time:   time.Now(),
		},
		StartupInfo: message.StartupInfo{
			DockerVersion: v.Version,
			Os:            v.Os,
			KernelVersion: v.KernelVersion,
			APIVersion:    v.APIVersion,
		},
	}

	for _, container := range containers {
		if cfg.ShowRunning {
			inspect, err := cli.ContainerInspect(ctx, container.ID)
			if err != nil {
				continue
			}

			creation, err := time.Parse(time.RFC3339Nano, inspect.Created)
			if err != nil {
				creation = time.Now()
			}
			cEvnt <- message.Event{
				MetaData: message.MetaData{
					Type:   "container",
					Action: "running",
					Time:   creation,
				},
				StartupInfo: message.StartupInfo{},
				Container: message.Container{
					ID:    container.ID,
					Name:  strings.TrimLeft(inspect.Name, "/"),
					Image: inspect.Config.Image,
				},
			}
		}

		log.Printf("Running %s %s\n", container.ID[:10], container.Image)
	}

	filters := filters.NewArgs()

	if cfg.ContainerEvents {
		filters.Add("type", events.ContainerEventType)
	}

	if cfg.VolumeEvents {
		filters.Add("type", events.VolumeEventType)
	}

	if cfg.NetworkEvents {
		filters.Add("type", events.NetworkEventType)
	}

	events, cliErrors := cli.Events(ctx, types.EventsOptions{
		Since:   "",
		Until:   "",
		Filters: filters,
	})

	errorCount := 0
	for {
		select {
		case event := <-events:
			errorCount = 0
			if event.Type != "container" {
				if _, ok := event.Actor.Attributes["container"]; !ok {
					log.Printf("Discar event: %+v\n", event)
					continue
				}
			}
			cEvnt <- parseEvent(event)

		case err := <-cliErrors:
			if err == io.EOF {
				break
			}

			log.Printf("Error while receiving events from Docker server: %s", err)
			errorCount++
			if errorCount > 10 {
				log.Printf("Maximum errors count. Exit")
				os.Exit(1)
			}

		case <-ctx.Done():
			return &ContextDoneError{}
		}
	}
}
func parseEvent(event events.Message) message.Event {
	e := message.Event{
		MetaData: message.MetaData{
			Type:   event.Type,
			Action: event.Action,
			Time:   time.Unix(event.Time, 0),
		},
	}

	switch event.Type {
	case "container":
		e.Container.Name = event.Actor.Attributes["name"]
		e.Container.Image = event.Actor.Attributes["image"]
		e.ContainerStatus.Signal = event.Actor.Attributes["signal"]
		e.ContainerStatus.ExitCode = event.Actor.Attributes["exitCode"]
		e.Container.ID = event.ID

	case "volume":
		e.Volume.Destination = event.Actor.Attributes["destination"]
		e.Volume.ID = event.Actor.ID
		e.Container.ID = event.Actor.Attributes["container"]

	case "network":
		e.Network.Name = event.Actor.Attributes["name"]
		e.Network.Type = event.Actor.Attributes["type"]
		e.Network.ID = event.Actor.ID
		e.Container.ID = event.Actor.Attributes["container"]
	}

	return e
}
