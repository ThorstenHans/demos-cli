package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var Version = "dev" // default if not set at build time

var rootCmd = &cobra.Command{
	Use:     "demos",
	Short:   "Fermyon & Akamai Load Testing Demos for KubeCon NA",
	Version: Version,
}

func init() {
	demos := demo.LoadAll()
	if err := demo.ValidateDemos(demos); err != nil {
		fmt.Printf("Provided demos are invalid: %s", err)
		os.Exit(1)
	}
	for _, d := range demos {
		runCmd.AddCommand(buildCommandFor(d))
	}
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(configureCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
