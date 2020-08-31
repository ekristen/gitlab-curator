package commands

import (
	"fmt"

	"github.com/ekristen/gitlab-curator/pkg/common"
	"github.com/urfave/cli/v2"
)

type command struct{}

func (w *command) Execute(c *cli.Context) error {
	fmt.Println(common.AppVersion)

	return nil
}

func init() {
	cmd := command{}

	cliCmd := &cli.Command{
		Name:   "version",
		Usage:  "print version",
		Action: cmd.Execute,
	}

	common.RegisterCommand(cliCmd)
}
