package cmd

import (
	"io"
	"os"
	"time"

	client "github.com/tsuru/autoscalev2/cmd/plugin/autoscalev2/stub" // TODO: use real client
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
	app.Commands = []*cli.Command{
		NewCmdTrigger(),
	}
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

func setupClient(c *cli.Context) error {
	// TODO: implement
	return nil
}

func getClient(c *cli.Context) (client.Client, error) {
	// TODO: implement
	return client.Client{}, nil
}

func newClient(c *cli.Context) (client.Client, error) {
	opts := client.ClientOptions{Timeout: c.Duration("timeout")}
	if rpaasURL := c.String("rpaas-url"); rpaasURL != "" {
		return client.NewClientWithOptions(rpaasURL, c.String("rpaas-user"), c.String("rpaas-password"), opts)
	}

	return client.NewClientThroughTsuruWithOptions(c.String("tsuru-target"), c.String("tsuru-token"), c.String("tsuru-service"), opts)
}
