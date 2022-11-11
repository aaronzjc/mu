package main

import (
	"log"
	"os"

	"github.com/aaronzjc/mu/internal"

	"github.com/urfave/cli"
)

var (
	appName = "mu-agent"
	usage   = "run mu-agent server"
	desc    = `mu-agent is the worker to craw pages, since some websites are blocked in china, distributed agents can be useful. `
	version = "7.0"
)

func main() {
	app := *cli.NewApp()
	app.Name = appName
	app.Usage = usage
	app.Description = desc
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config,c",
			Usage: "(config) Load configuration from `FILE`",
		},
	}
	app.Before = internal.SetupAgent
	app.Action = internal.RunAgent

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
