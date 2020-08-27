package webhookdocker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/logocomune/webhookdocker/container"
	"github.com/logocomune/webhookdocker/message"
	"github.com/logocomune/webhookdocker/processor"
	"github.com/logocomune/webhookdocker/webhook"
)

const (
	httpClientTimeOut = 3 * time.Second
	aggregationTime   = 3 * time.Second
)

type sender interface {
	Send(events map[string]message.ContainerEventsGroup)
}

func MainProcess(cfg CommonCfg, kb Keybase, sl Slack, ex WebEx, appVersion, builtDate string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	defer func() {
		signal.Stop(shutdown)
		cancel()
	}()

	go func() {
		select {
		case <-shutdown:
			log.Println("Shutdown request")
			cancel()
		case <-ctx.Done():
		}
	}()

	nodeName := getNodeName(cfg)

	formatter, err := message.New(nodeName)
	if err != nil {
		log.Fatalln("formatter error", err)
	}

	var wbSender sender

	if kb.Endpoint != "" {
		wbSender = webhook.NewKB(kb.Endpoint, httpClientTimeOut, cfg.Docker.ExternalInstanceInspection)
	}

	if sl.Endpoint != "" {
		wbSender = webhook.NewSlack(sl.Endpoint, httpClientTimeOut, cfg.Docker.ExternalInstanceInspection)
	}

	if ex.Endpoint != "" {
		wbSender = webhook.NewWebEx(ex.Endpoint, httpClientTimeOut, cfg.Docker.ExternalInstanceInspection)
	}

	if wbSender == nil {
		log.Fatalln("No endpoint chosen")
	}

	processor := processor.NewProcessor(ctx, aggregationTime, formatter, wbSender)

	return container.DockerEvents(ctx, processor.Q, container.DockerCfg{
		ContainerEvents:   cfg.Docker.Listen.ContainerEvents,
		VolumeEvents:      cfg.Docker.Listen.VolumeEvents,
		NetworkEvents:     cfg.Docker.Listen.ContainerEvents,
		ShowRunning:       cfg.Docker.ShowRunning,
		ContainerActions:  cfg.Docker.Listen.ContainerActions,
		NetworkActions:    cfg.Docker.Listen.NetworkActions,
		VolumeActions:     cfg.Docker.Listen.VolumeActions,
		FilterName:        cfg.Docker.Filter.ContainerName,
		NegateFilterName:  cfg.Docker.Filter.NegateContainerName,
		FilterImage:       cfg.Docker.Filter.ImageName,
		NegateFilterImage: cfg.Docker.Filter.NegateImageName,
	}, appVersion, builtDate)
}

func getNodeName(cfg CommonCfg) string {
	nodeName := ""
	if !cfg.HideNodeName {
		nodeName = cfg.NodeName

		if nodeName == "" {
			nodeName, _ = os.Hostname()
		}
	}

	return nodeName
}
