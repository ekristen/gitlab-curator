package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// GlobalFlags --
func GlobalFlags() []cli.Flag {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log Level",
			EnvVars: []string{"LOG_LEVEL"},
			Value:   "info",
		},
	}

	return globalFlags
}

// GlobalBefore --
func GlobalBefore(c *cli.Context) error {
	switch c.String("log-level") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	return nil
}
