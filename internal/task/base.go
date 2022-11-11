package task

import "github.com/urfave/cli"

type Task interface {
	Name() string
	Run(*cli.Context) error
}
