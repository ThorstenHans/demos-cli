package commands

import "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Run a particular demo",
	Long: `Run a particular demo

These sub-commands are dynamically loaded upon application start. 
All command names and aliases are generate from your demo definition. 
Your demos are stored in $HOME/.demos/demos.json
`,
}
