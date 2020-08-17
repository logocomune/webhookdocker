package webhook

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/logocomune/webhookdocker/message"
)

const dockerWebhook = "docker-webhook"

type KB struct {
	webHookUlr string
	client     *http.Client
	formatter  formatter
}

//NewKB Initialize Keybase webhook sender
func NewKB(webHookUlr string, timeOut time.Duration, externalInspectUrl string) *KB {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeOut,
		}).DialContext,
	}

	return &KB{
		webHookUlr: webHookUlr,
		client: &http.Client{
			Transport: tr,
			Timeout:   timeOut,
		},
		formatter: formatter{
			labels:             keybaseLabels,
			externalInspectUrl: externalInspectUrl},
	}
}

//Send Send a group of docker events to keybase webhook
func (k *KB) Send(events map[string]message.ContainerEventsGroup) {
	msg := ""

	if e, ok := events[dockerWebhook]; ok {
		str, _ := eventsToStr(k.formatter, e.NodeName, e)
		msg += str + "\n\n"

		delete(events, dockerWebhook)
	}

	for _, g := range events {
		str, _ := eventsToStr(k.formatter, g.NodeName, g)
		msg += str + "\n\n"
	}

	if msg == "" {
		return
	}

	if err := postWebHook(k.client, k.webHookUlr, struct{ Msg string }{Msg: msg}); err != nil {
		log.Println("Keybase, POST error:", err.Error())
	}
}
