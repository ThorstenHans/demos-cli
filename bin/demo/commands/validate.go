package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates your demos",
	Run: func(cmd *cobra.Command, args []string) {
		if demo.HasDemosFile() {
			demos := demo.LoadAll()
			if err := demo.ValidateDemos(demos); err != nil {
				fmt.Printf("Demos not valid :( (%s)\n", err)
				os.Exit(1)
			}
			fmt.Println("All demos are valid! Let's go!! 🚀")
			os.Exit(0)
		}
		fmt.Printf("There are no custom demos yet!\nConsider running `demos eject`\n")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
