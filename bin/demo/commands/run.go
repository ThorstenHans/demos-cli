package commands

import "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Run Demos",
}
