package commands

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/thorstenhans/demos-over-ssh/internal/demo"
)

var (
	username string
	port     int
	ip       net.IP
)

var setCmd = &cobra.Command{
	Use:     "set-config",
	Short:   "Set configuration values",
	GroupID: "config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return demo.Configure()
	},
}

var getCmd = &cobra.Command{
	Use:     "get-config",
	Short:   "Get configuration values",
	GroupID: "config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := demo.LoadConfig()
		if err != nil {
			fmt.Printf("Failed to retrieve the configuration")
			os.Exit(1)
		}
		fmt.Println("Current Configuration:")
		fmt.Printf(" Jump Box: %s\n", cfg.GetJumpBoxEndpoint())
		fmt.Printf(" Username: %s\n", cfg.JumpBoxUser)
		return nil
	},
}

func init() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "config",
		Title: "config",
	})
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	//rootCmd.AddCommand(configCmd)
}
