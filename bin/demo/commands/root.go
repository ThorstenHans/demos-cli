package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var Version = "dev" // default if not set at build time

var rootCmd = &cobra.Command{
	Use:   "demos",
	Short: "Demo or Die! Run demos over SSH",
	Long: `
Demo 🎬 or Die 💀! 
	
The demos CLI allows running even complex demos over SSH. 
Meaning only stdin and stdout must be transferred over the network. 
This is especially helpful to overcome bad network connectivity.

Enjoy the conference 🤘`,
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
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

}

func Execute() error {
	return rootCmd.Execute()
}
