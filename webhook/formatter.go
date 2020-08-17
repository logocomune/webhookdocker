package webhook

import (
	"strings"
	"time"

	"github.com/logocomune/webhookdocker/message"
)

const (
	appTitle = iota
	startupTitle1
	startupTitle2
	groupTitle
	groupFooter
	groupFooterWithInspect
	containerDefault
	containerExit0
	containerDie
	containerKill
	volumeMount
	volumeUnmount
	networkDefault

	idMaxSize = 12

	githubUrl    = "https://github.com/logocomune/webhookdocker/issues"
	dockerHubUrl = "https://hub.docker.com/r/logocomune/webhook-docker"
)

type formatter struct {
	labels             map[int]string
	externalInspectUrl string
}

func (f formatter) startupMessage(nodeName string, info message.StartupInfo) string {
	if info.DockerVersion == "" {
		return ""
	}

	msg := f.labels[appTitle] + f.labels[startupTitle1] + f.labels[startupTitle2]

	msg = strings.Replace(msg, "__APP_VERSION__", info.AppVersion, -1)
	msg = strings.Replace(msg, "__APP_BUILT_DATE__", info.AppBuiltDate, -1)

	msg = strings.Replace(msg, "__GITHUB_URL__", githubUrl, -1)
	msg = strings.Replace(msg, "__DOCKER_HUB_URL__", dockerHubUrl, -1)

	msg = strings.Replace(msg, "__DOCKER_VERSION__", info.DockerVersion, -1)
	msg = strings.Replace(msg, "__DOCKER_API_VERSION__", info.APIVersion, -1)
	msg = strings.Replace(msg, "__OS__", info.Os, -1)
	msg = strings.Replace(msg, "__KERNEL_VERSION__", info.KernelVersion, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	msg = strings.Replace(msg, "__TAB__", "\t", -1)
	hostnameReplacer := ""

	if nodeName != "" {
		hostnameReplacer = nodeName
	}

	msg = strings.Replace(msg, "__NODE_NAME__", hostnameReplacer, -1)

	return msg
}

func (f formatter) titleMessage(name string, image string, nodeName string, t time.Time) string {
	time := t.Format(time.RFC3339)
	msg := f.labels[groupTitle]
	msg = strings.Replace(msg, "__IMAGE__", image, -1)
	msg = strings.Replace(msg, "__NAME__", name, -1)
	msg = strings.Replace(msg, "__TIME__", time, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	hostnameReplacer := ""

	if nodeName != "" {
		hostnameReplacer = "*@*_" + nodeName + "_ "
	}

	msg = strings.Replace(msg, "__NODE_NAME__", hostnameReplacer, -1)
	msg = strings.Replace(msg, "__NODE_NAME_W_SPACES__", strings.Replace(hostnameReplacer, "*@*_", " *@* _", 1), -1)

	return msg
}

func (f formatter) footerMessage(id string) string {
	msg := f.labels[groupFooter]
	var inspectUrl string
	if f.externalInspectUrl != "" {
		msg = f.labels[groupFooterWithInspect]
		inspectUrl = buildInspectUrl(f.externalInspectUrl, id)
	}

	id = shortId(id)
	msg = strings.Replace(msg, "__ID__", id, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	msg = strings.Replace(msg, "__INSPECT_URL__", inspectUrl, -1)

	return msg
}

func buildInspectUrl(url string, id string) string {
	url = strings.Replace(url, "__ID__", id, -1)
	url = strings.Replace(url, "__SHORT_ID__", shortId(id), -1)
	return url
}

func shortId(id string) string {
	if len(id) <= idMaxSize {
		return id
	}
	return id[:idMaxSize]
}

func (f formatter) containerMessage(meta message.MetaData, eContainer message.Container, eStatus message.ContainerStatus) string {
	msg := ""

	switch meta.Action {
	case "kill":
		msg = f.labels[containerKill]

	case "die":
		msg = f.labels[containerDie]
		if eStatus.ExitCode == "0" {
			msg = f.labels[containerExit0]
		}

	default:
		msg = f.labels[containerDefault]
	}

	instanceID := eContainer.ID
	action := meta.Action
	image := eContainer.Image
	name := eContainer.Name
	time := meta.Time.Format(time.RFC3339)
	exitCode := eStatus.ExitCode
	signal := eStatus.Signal

	msg = strings.Replace(msg, "__ID__", instanceID, -1)
	msg = strings.Replace(msg, "__ACTION__", strings.Title(action), -1)
	msg = strings.Replace(msg, "__IMAGE__", image, -1)
	msg = strings.Replace(msg, "__NAME__", name, -1)
	msg = strings.Replace(msg, "__TIME__", time, -1)
	msg = strings.Replace(msg, "__EXIT_CODE__", exitCode, -1)
	msg = strings.Replace(msg, "__SIGNAL__", signal, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	msg = strings.Replace(msg, "__TAB__", "\t", -1)

	return msg
}

func (f formatter) volumeMessage(meta message.MetaData, volume message.Volume) string {
	id := volume.ID
	dest := volume.Destination
	action := meta.Action

	msg := ""

	switch action {
	case "mount":
		msg = f.labels[volumeMount]

	case "unmount":
		msg = f.labels[volumeUnmount]
	}

	msg = strings.Replace(msg, "__ACTION__", strings.Title(action), -1)
	msg = strings.Replace(msg, "__VOLUME_ID__", id, -1)
	msg = strings.Replace(msg, "__VOLUME_DESTINATION__", dest, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	msg = strings.Replace(msg, "__TAB__", "\t", -1)

	return msg
}

func (f formatter) networkMessage(meta message.MetaData, network message.Network) string {
	msg := f.labels[networkDefault]
	name := network.Name
	id := network.ID
	action := meta.Action

	msg = strings.Replace(msg, "__ACTION__", strings.Title(action), -1)
	msg = strings.Replace(msg, "__NETWORK_ID__", id, -1)
	msg = strings.Replace(msg, "__NETWORK_NAME__", name, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	msg = strings.Replace(msg, "__TAB__", "\t", -1)

	return msg
}

func buildMsg(events map[string]message.ContainerEventsGroup, w *WebEx) (string, bool) {
	msg := ""

	if e, ok := events[dockerWebhook]; ok {
		str, _ := eventsToStr(w.formatter, e.NodeName, e)
		msg += str + "\n\n"

		delete(events, dockerWebhook)
	}

	for _, g := range events {
		str, _ := eventsToStr(w.formatter, g.NodeName, g)
		msg += str + "\n\n"
	}

	if msg == "" {
		return "", false
	}

	return msg, true
}

//eventsToStr Create a message string for a group of container events
func eventsToStr(f formatter, nodeName string, eventsGroup message.ContainerEventsGroup) (string, bool) {
	s := f.startupMessage(nodeName, eventsGroup.StartupInfo)
	if len(eventsGroup.Events) == 0 && s == "" {
		return "", false
	}

	nEvents := len(eventsGroup.Events)

	containerID := eventsGroup.ID
	containerImage := eventsGroup.Image

	containerName := eventsGroup.Name

	for idx, event := range eventsGroup.Events {
		if idx == 0 {
			s += f.titleMessage(containerName, containerImage, nodeName, event.Time)
		}

		switch event.MetaData.Type {
		case "container":
			s += f.containerMessage(event.MetaData, event.Container, event.ContainerStatus)
		case "volume":
			s += f.volumeMessage(event.MetaData, event.Volume)

		case "network":
			s += f.networkMessage(event.MetaData, event.Network)
		}

		if idx != nEvents-1 {
			s += "> \n"
		} else {
			s += f.footerMessage(containerID)
		}
	}

	return s, true
}
