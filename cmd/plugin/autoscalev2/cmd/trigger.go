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
			NewCmdTriggerAdd(),
			// NewCmdTriggerDelete(),
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
				Aliases:  []string{"i"},
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

func NewCmdTriggerAdd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a trigger on the instance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "instance",
				Aliases:  []string{"i"},
				Usage:    "the reverse proxy instance name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name for the trigger",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "Type of the trigger (https://keda.sh/docs/scalers)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"metadata", "c"},
				Usage:    "Configuration for the trigger (metadata schema depends on the trigger type)",
				Required: true,
			},
		},
		Before: setupClient,
		Action: runAddTrigger,
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

func runAddTrigger(c *cli.Context) error {
	cli, err := getClient(c)
	if err != nil {
		return err
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(c.String("config")), &config); err != nil {
		return fmt.Errorf("Config could not be parsed. Not a valid JSON: %w", err)
	}

	args := client.UpsertTriggerArgs{
		Instance: c.String("instance"),
		Name:     c.String("name"),
		Type:     c.String("type"),
		Metadata: config,
	}

	if err := cli.UpsertTrigger(c.Context, args); err != nil {
		return fmt.Errorf("Error creating/updating trigger %q: %w", c.String("name"), err)
	}

	return nil
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
