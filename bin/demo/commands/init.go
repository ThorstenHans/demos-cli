package commands

import (
	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var initCmd = &cobra.Command{
	Use:     "initialize",
	Aliases: []string{"init", "eject"},
	Short:   "Initialize the demo file",
	Long: `Initialize the demo file

The demo file (demos.json) is written to %HOME/demos/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return demo.GenerateSampleDemosFile()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
