package webhook

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/logocomune/webhookdocker/message"
)

type WebEx struct {
	webHookUlr string
	formatter  formatter
	client     *http.Client
}

type webexMessage struct {
	Markdown string `json:"markdown"`
}

//NewWebEx Initialize WebEx webhook sender
func NewWebEx(webHookUlr string, timeOut time.Duration, externalInspectUrl string) *WebEx {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeOut,
		}).DialContext,
	}

	return &WebEx{
		webHookUlr: webHookUlr,
		client: &http.Client{
			Transport: tr,
			Timeout:   timeOut,
		},
		formatter: formatter{labels: webExLabels,
			externalInspectUrl: externalInspectUrl},
	}
}

// WebEx doc
// https://apphub.webex.com/teams/applications/incoming-webhooks-cisco-systems
func (w *WebEx) Send(events map[string]message.ContainerEventsGroup) {
	msg, ok := buildMsg(events, w)
	if !ok {
		return
	}

	wMsg := webexMessage{
		Markdown: msg,
	}

	if err := postWebHook(w.client, w.webHookUlr, wMsg); err != nil {
		log.Println("WebEx, POST error:", err.Error())
	}
}
