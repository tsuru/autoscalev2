package cmd

import (
	"io"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func NewDefaultApp() *cli.App {
	return NewApp(os.Stdout, os.Stderr)
}

func NewApp(o, e io.Writer) (app *cli.App) {
	app = cli.NewApp()
	app.Usage = "Manages Autoscalev2 Instances"
	// app.Version = version.Version
	app.ErrWriter = e
	app.Writer = o
	app.Commands = []*cli.Command{}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "tsuru-target",
			Usage:   "address of Tsuru server",
			EnvVars: []string{"TSURU_TARGET"},
		},
		&cli.StringFlag{
			Name:        "tsuru-token",
			Usage:       "authentication credential to Tsuru server",
			EnvVars:     []string{"TSURU_TOKEN"},
			DefaultText: "-",
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "time limit that a remote operation (HTTP request) can take",
			Value: 60 * time.Second,
		},
	}
	return app
}
