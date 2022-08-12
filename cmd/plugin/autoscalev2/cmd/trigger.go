package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	client "github.com/tsuru/autoscalev2/cmd/plugin/autoscalev2/stub"
	clientTypes "github.com/tsuru/autoscalev2/cmd/plugin/autoscalev2/stub"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func NewCmdTrigger() *cli.Command {
	return &cli.Command{
		Name:  "triggers",
		Usage: "Manages autoscalev2 triggers",
		Subcommands: []*cli.Command{
			NewCmdTriggerList(),
			NewCmdTriggerAdd(),
			NewCmdTriggerDelete(),
			NewCmdTriggerGet(),
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

func runAddTrigger(c *cli.Context) error {
	cli, err := getClient(c)
	if err != nil {
		return err
	}

	var config clientTypes.TriggerMetadata
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

func NewCmdTriggerDelete() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a trigger on the instance",
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
		},
		Before: setupClient,
		Action: runDeleteTrigger,
	}
}

func runDeleteTrigger(c *cli.Context) error {
	cli, err := getClient(c)
	if err != nil {
		return err
	}

	args := client.DeleteTriggerArgs{
		Instance: c.String("instance"),
		Name:     c.String("name"),
	}

	if err := cli.DeleteTrigger(c.Context, args); err != nil {
		return fmt.Errorf("Error deleting trigger %q: %w", c.String("name"), err)
	}

	return nil
}

func NewCmdTriggerGet() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get configuration for a trigger on the instance",
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
		},
		Before: setupClient,
		Action: runGetTrigger,
	}
}

func runGetTrigger(c *cli.Context) error {
	cli, err := getClient(c)
	if err != nil {
		return err
	}

	args := client.GetTriggerArgs{
		Instance: c.String("instance"),
		Name:     c.String("name"),
	}

	trigger, err := cli.GetTrigger(c.Context, args)
	if err != nil {
		return fmt.Errorf("Error getting trigger %q: %w", c.String("name"), err)
	}

	var buf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2)
	printMetadata := map[string]interface{}{
		"Configuration": trigger.Metadata,
	}
	yamlEncoder.Encode(printMetadata)

	fmt.Fprintf(c.App.Writer, "Instance: %s\n", args.Instance)
	fmt.Fprintf(c.App.Writer, "Trigger Name: %s\n", trigger.Name)
	fmt.Fprintf(c.App.Writer, "Trigger Type: %s\n", trigger.Type)
	fmt.Fprintln(c.App.Writer, "")
	fmt.Fprintln(c.App.Writer, buf.String())

	return nil
}
