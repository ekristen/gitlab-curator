package main

import (
	"os"
	"path"

	"github.com/ekristen/gitlab-curator/pkg/common"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/ekristen/gitlab-curator/pkg/commands/run"
	_ "github.com/ekristen/gitlab-curator/pkg/commands/version"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "gitlab-curator"
	app.Version = common.AppVersion.Summary
	app.Authors = []*cli.Author{
		{
			Name:  "Erik Kristensen",
			Email: "erik@erikkristensen.com",
		},
	}

	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Fatalln("Command", command, "not found.")
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}
