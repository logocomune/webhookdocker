package webhookdocker

type CommonCfg struct {
	NodeName     string
	HideNodeName bool `conf:"default:false"`
	Docker       struct {
		ShowRunning bool `conf:"default:false"`
		Filter      struct {
			ContainerName string `conf:""`
			ImageName     string `conf:""`
		}
		Listen struct {
			ContainerEvents  bool     `conf:"default:true"`
			NetworkEvents    bool     `conf:"default:true"`
			VolumeEvents     bool     `conf:"default:true"`
			ContainerActions []string `conf:"default:attach;create;destroy;detach;die;kill;oom;pause;rename;restart;start;stop;unpause;update"`
			NetworkActions   []string `conf:"default:create;connect;destroy;disconnect;remove"`
			VolumeActions    []string `conf:"default:create;destroy;mount;unmount"`
		}
	}
}

type Keybase struct {
	Endpoint string `conf:""`
}

type Slack struct {
	Endpoint string `conf:""`
}

type WebEx struct {
	Endpoint string `conf:"flag:webex-endpoint"`
}
