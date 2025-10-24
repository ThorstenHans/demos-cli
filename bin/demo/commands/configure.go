package commands

import (
	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "(Re)Configure the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return demo.Configure()
	},
}
