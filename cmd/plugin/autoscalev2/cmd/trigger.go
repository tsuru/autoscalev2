package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	client "github.com/tsuru/autoscalev2/cmd/plugin/autoscalev2/stub"
	clientTypes "github.com/tsuru/autoscalev2/cmd/plugin/autoscalev2/stub"
	"github.com/urfave/cli/v2"
)

func NewCmdTrigger() *cli.Command {
	return &cli.Command{
		Name:  "triggers",
		Usage: "Manages autoscalev2 triggers",
		Subcommands: []*cli.Command{
			NewCmdTriggerList(),
		},
	}
}

func NewCmdTriggerList() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "Shows the triggers on the instance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "instance",
				Aliases:  []string{"tsuru-service-instance", "i"},
				Usage:    "the reverse proxy instance name",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "raw-output",
				Aliases: []string{"r"},
				Usage:   "show as JSON instead of table format",
				Value:   false,
			},
		},
		Before: setupClient,
		Action: runListTriggers,
	}
}

func runListTriggers(c *cli.Context) error {
	cli, err := getClient(c)
	if err != nil {
		return err
	}

	args := client.ListTriggersArgs{Instance: c.String("instance")}
	routes, err := cli.ListTriggers(c.Context, args)
	if err != nil {
		return err
	}

	if c.Bool("raw-output") {
		return writeTriggersListRawOutput(c.App.Writer, routes)
	}
	return writeTriggersListSimple(c.App.Writer, routes)
}

func writeTriggersListSimple(w io.Writer, triggers []clientTypes.Trigger) error {
	data := [][]string{}
	for _, t := range triggers {
		data = append(data, []string{t.Name, t.Type})
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Type"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
	table.AppendBulk(data)
	table.Render()
	return nil
}

func writeTriggersListRawOutput(w io.Writer, triggers []clientTypes.Trigger) error {
	message, err := json.MarshalIndent(triggers, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(message))
	return nil
}
