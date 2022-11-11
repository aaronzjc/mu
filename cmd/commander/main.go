package main

import (
	"log"
	"os"

	"github.com/aaronzjc/mu/internal"

	"github.com/urfave/cli"
)

var (
	appName = "mu-commander"
	usage   = "run mu-commander server"
	desc    = `mu-commander is the schedule server for dispatching craw jobs, monitor agent status, etc. `
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
	app.Before = internal.SetupCommander
	app.Action = internal.RunCommander

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
