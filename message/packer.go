package message

import (
	lru "github.com/hashicorp/golang-lru"
)

const (
	cacheSize     = 128
	dockerWebhook = "docker-webhook"
)

type Packer struct {
	cache    *lru.Cache
	nodeName string
}

//New Initialize new formatter
func New(nodeName string) (*Packer, error) {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	return &Packer{
		cache:    cache,
		nodeName: nodeName,
	}, nil
}

//EventPacker Receive a list of docker events and group them by container id
func (p *Packer) EventPacker(events []Event) map[string]ContainerEventsGroup {
	aggrMsgs := make(map[string]ContainerEventsGroup)

	for _, e := range events {
		if e.MetaData.Type == dockerWebhook {
			id := dockerWebhook
			aggrMsgs[id] = ContainerEventsGroup{
				Group: Group{
					NodeName: p.nodeName,
				},
				StartupInfo: e.StartupInfo,
			}

			continue
		}

		id := e.Container.ID

		eventsGroup := aggrMsgs[id]
		eventsGroup.Events = append(eventsGroup.Events, e)

		if e.MetaData.Type == "container" {
			eventsGroup.Container = e.Container
			p.cache.Add(id, e.Container)
		}

		if eventsGroup.Container.Name == "" {
			if c, ok := p.cache.Get(id); ok {
				container := c.(Container)
				eventsGroup.Container = container
			}
		}

		eventsGroup.NodeName = p.nodeName
		aggrMsgs[id] = eventsGroup
	}

	for k := range aggrMsgs {
		if k != dockerWebhook && aggrMsgs[k].ID == "" {
			delete(aggrMsgs, k)
		}
	}

	return aggrMsgs
}
