package webhook

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/logocomune/webhookdocker/message"
)

type Slack struct {
	webHookUlr string
	formatter  formatter
	client     *http.Client
}

type slackMessage struct {
	Blocks []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string    `json:"type"`
	Text slackText `json:"text"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewSlack(webHookUlr string, timeOut time.Duration, externalInspectUrl string) *Slack {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeOut,
		}).DialContext,
	}

	return &Slack{
		webHookUlr: webHookUlr,
		client: &http.Client{
			Transport: tr,
			Timeout:   timeOut,
		},
		formatter: formatter{
			labels:             slackLabels,
			externalInspectUrl: externalInspectUrl,
		},
	}
}

func (s *Slack) Send(events map[string]message.ContainerEventsGroup) {
	msg := ""

	if e, ok := events[dockerWebhook]; ok {
		str, _ := eventsToStr(s.formatter, e.NodeName, e, ">")
		msg += str + "\n\n"

		delete(events, dockerWebhook)
	}

	for _, g := range events {
		str, _ := eventsToStr(s.formatter, g.NodeName, g, ">")
		msg += str + "\n\n"
	}

	if msg == "" {
		return
	}

	sMsg := slackMessage{
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: slackText{
					Type: "mrkdwn",
					Text: msg,
				},
			},
		},
	}

	if err := postWebHook(s.client, s.webHookUlr, sMsg); err != nil {
		log.Println("Slack, POST error:", err.Error())
	}
}
