package main

import (
	"fmt"
	"os"

	"github.com/tsuru/autoscalev2/cmd/plugin/autoscalev2/cmd"
)

func main() {
	app := cmd.NewDefaultApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(app.ErrWriter, err)
		os.Exit(1)
	}
}
