package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/logocomune/webhookdocker"
	"github.com/logocomune/webhookdocker/container"
	"github.com/pkg/errors"
)

var (
	//version is the application version. (Injected by make)
	version string
	//shortVersion is the application version. It's used also in frontend and api response (Injected by make)
	shortVersion = "v0.0.0"
	//commit commit hash. (Injected by make)
	commit string
	//branch current branch name. (Injected by make)
	branch string
	//buildDate Build date. (Injected by make)
	buildDate string

	build = "develop"
)

type cfgArgs struct {
	webhookdocker.CommonCfg
	Keybase webhookdocker.Keybase
	Slack   webhookdocker.Slack
	WebEx   webhookdocker.WebEx
}

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	if err := run(); err != nil {
		var cDoneErro *container.ContextDoneError
		if errors.As(err, &cDoneErro) {
			os.Exit(0)
		}

		log.Println("error :", err)

		os.Exit(1)
	}
}

func run() error {
	var cfg cfgArgs

	log.Printf("main : Started : Application initializing : version %s (Built: %s)", shortVersion, buildDate)
	log.Println("Repository: https://github.com/logocomune/webhookdocker")
	log.Println(os.Args[0], "-h", "for help")

	if err := conf.Parse(os.Args[1:], "WD", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("WD", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}

			fmt.Println(usage)

			return nil
		}

		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}

	log.Printf("main : Config :\n%v\n", out)

	defer log.Println("main : Completed")

	return webhookdocker.MainProcess(cfg.CommonCfg, cfg.Keybase, cfg.Slack, cfg.WebEx, shortVersion, buildDate)
}
