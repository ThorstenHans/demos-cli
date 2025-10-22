package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

func buildCommandFor(d demo.DemoScript) *cobra.Command {

	return &cobra.Command{
		Use:     d.Command,
		Aliases: []string{d.Alias},
		Short:   d.ShortDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := demo.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading Configuration: %s\n", err)
				os.Exit(1)
			}
			return demo.Run(d, cfg)
		},
	}
}
