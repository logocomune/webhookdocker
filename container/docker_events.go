package container

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/logocomune/webhookdocker/message"
)

type DockerCfg struct {
	ContainerEvents   bool
	VolumeEvents      bool
	NetworkEvents     bool
	ShowRunning       bool
	ContainerActions  []string
	NetworkActions    []string
	VolumeActions     []string
	FilterName        string
	NegateFilterName  bool
	FilterImage       string
	NegateFilterImage bool
}

const maxErrors = 10

type ContextDoneError struct {
}

func (c *ContextDoneError) Error() string { return "DockerEvents: Context done. Exit" }

func DockerEvents(ctx context.Context, cEvnt chan message.Event, cfg DockerCfg, appVersion, appBuiltDate string) error {
	eventFilter, err := newFilter(cfg)
	if err != nil {
		return err
	}

	//cli, err := docker.NewEnvClient()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})

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
			AppVersion:    appVersion,
			AppBuiltDate:  appBuiltDate,
		},
	}

	for _, c := range containers {
		if cfg.ShowRunning {
			inspect, err := cli.ContainerInspect(ctx, c.ID)
			if err != nil {
				continue
			}

			creation, err := time.Parse(time.RFC3339Nano, inspect.Created)
			if err != nil {
				creation = time.Now()
			}
			cEvnt <- message.Event{
				MetaData: message.MetaData{
					Type:   "c",
					Action: "running",
					Time:   creation,
				},
				StartupInfo: message.StartupInfo{},
				Container: message.Container{
					ID:    c.ID,
					Name:  strings.TrimLeft(inspect.Name, "/"),
					Image: inspect.Config.Image,
				},
			}
		}

		log.Printf("Running %s %s\n", c.ID[:10], c.Image)
	}

	filterItems := filters.NewArgs()

	if cfg.ContainerEvents {
		filterItems.Add("type", string(events.ContainerEventType))
	}

	if cfg.VolumeEvents {
		filterItems.Add("type", string(events.VolumeEventType))
	}

	if cfg.NetworkEvents {
		filterItems.Add("type", string(events.NetworkEventType))
	}

	evnts, cliErrors := cli.Events(ctx, events.ListOptions{
		Since:   "",
		Until:   "",
		Filters: filterItems,
	})

	errorCount := 0

	for {
		select {
		case event := <-evnts:
			errorCount = 0

			if !eventFilter.accept(event) {
				continue
			}

			if event.Type != events.ContainerEventType {
				if _, ok := event.Actor.Attributes["c"]; !ok {
					continue
				}
			}
			cEvnt <- parseEvent(event)

		case err := <-cliErrors:
			if err == io.EOF {
				break
			}

			log.Printf("Error while receiving evnts from Docker server: %s", err)
			errorCount++

			if errorCount > maxErrors {
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
			Type:   string(event.Type),
			Action: string(event.Action),
			Time:   time.Unix(event.Time, 0),
		},
	}

	switch event.Type {
	case "container":
		e.Container.Name = event.Actor.Attributes["name"]
		e.Container.Image = event.Actor.Attributes["image"]
		e.ContainerStatus.Signal = event.Actor.Attributes["signal"]
		e.ContainerStatus.ExitCode = event.Actor.Attributes["exitCode"]
		e.Container.ID = event.Actor.ID

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
