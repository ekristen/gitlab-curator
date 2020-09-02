package run

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"

	"github.com/ekristen/gitlab-curator/pkg/commands"
	"github.com/ekristen/gitlab-curator/pkg/common"
	"github.com/ekristen/gitlab-curator/pkg/types"
)

type command struct{}

func (s *command) Execute(c *cli.Context) error {
	token := c.String("token")
	sourceType := c.String("source-type")
	sourceID := c.String("source-id")
	file := c.String("file")
	dryrun := c.Bool("dry-run")

	if !fileExists(file) {
		return fmt.Errorf("%s not found or not readable", file)
	}

	logrus.WithFields(logrus.Fields{
		"token":      token,
		"sourceType": sourceType,
		"sourceID":   sourceID,
		"file":       file,
	}).Debug("info")

	git, err := gitlab.NewClient(token)
	if err != nil {
		panic(err)
	}

	options := types.NewOptions(git, sourceType, sourceID, dryrun)

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var policy types.Policy

	err = yaml.Unmarshal(yamlFile, &policy)
	if err != nil {
		return err
	}

	if err := policy.ResourceRules.Process(options); err != nil {
		return err
	}

	return nil
}

func init() {
	cmd := command{}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Usage:    "GitLab Token",
			EnvVars:  []string{"TOKEN"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "source-type",
			Usage:   "Source Type",
			EnvVars: []string{"SOURCE_TYPE"},
			Value:   "group",
		},
		&cli.StringFlag{
			Name:     "source-id",
			Usage:    "Source ID",
			EnvVars:  []string{"SOURCE_ID"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "file",
			Usage:    "File",
			EnvVars:  []string{"FILE"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Usage:   "Dry Run",
			EnvVars: []string{"DRY_RUN"},
			Value:   false,
		},
	}

	cliCmd := &cli.Command{
		Name:   "run",
		Usage:  "run a policy file",
		Action: cmd.Execute,
		Flags:  append(flags, commands.GlobalFlags()...),
		Before: commands.GlobalBefore,
	}

	common.RegisterCommand(cliCmd)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
