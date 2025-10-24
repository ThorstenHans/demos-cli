package commands

import (
	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var ejectCmd = &cobra.Command{
	Use:   "eject",
	Short: "Write default demos to demo.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		return demo.GenerateSampleDemosFile()
	},
}

func init() {
	rootCmd.AddCommand(ejectCmd)
}
