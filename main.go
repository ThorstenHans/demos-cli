package main

import (
	"fmt"
	"os"

	"github.com/thorstenhans/demos-over-ssh/bin/demo/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
