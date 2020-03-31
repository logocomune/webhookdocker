package processor

import (
	"context"
	"log"
	"time"

	"github.com/logocomune/webhookdocker/message"
)

const (
	queueSize = 40
)

type Processor struct {
	Q         chan message.Event
	sender    sender
	formatter formatter
}

type sender interface {
	Send(events map[string]message.ContainerEventsGroup)
}

type formatter interface {
	EventPacker([]message.Event) map[string]message.ContainerEventsGroup
}

func NewProcessor(ctx context.Context, d time.Duration, formatter formatter, sender sender) *Processor {
	qtardis := Processor{
		sender:    sender,
		formatter: formatter,
		Q:         make(chan message.Event, queueSize),
	}
	qtardis.startUp(ctx, d)

	return &qtardis
}

func (q *Processor) startUp(ctx context.Context, d time.Duration) {
	buffer := make([]message.Event, 0, queueSize)
	ticker := time.NewTicker(d)
	working := false

	go func() {
		for {
			select {
			case <-ticker.C:
				if len(buffer) == 0 {
					continue
				}

				if working {
					log.Println("Sending in progress... skip")
					continue
				}

				working = true
				//Execute same action that put working in false
				groupOfEvents := q.formatter.EventPacker(buffer)

				go func() {
					q.sender.Send(groupOfEvents)

					working = false
				}()

				//Init Q
				buffer = make([]message.Event, 0, queueSize)

			case m, ok := <-q.Q:
				if !ok {
					continue
				}

				buffer = append(buffer, m)

			case <-ctx.Done():
				ticker.Stop()
				close(q.Q)

				return
			}
		}
	}()
}
