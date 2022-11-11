package main

import (
	"log"
	"os"

	"github.com/aaronzjc/mu/internal"

	"github.com/urfave/cli"
)

var (
	appName = "mu-api"
	usage   = "run mu-api server"
	desc    = `mu-api is the api server for mu project, it provides index api & admin api services. `
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
	app.Before = internal.SetupApi
	app.Action = internal.RunApi

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
