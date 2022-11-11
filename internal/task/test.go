package task

import (
	"fmt"

	"github.com/urfave/cli"
)

type Test struct{}

func (t *Test) Name() string {
	return "echo"
}

func (t *Test) Run(ctx *cli.Context) error {
	fmt.Println("hello world")
	return nil
}

var test Test

func NewTestTask() cli.Command {
	return cli.Command{
		Name:   test.Name(),
		Usage:  test.Name() + " run",
		Action: test.Run,
	}
}
