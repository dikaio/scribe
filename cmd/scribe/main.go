package main

import (
	"fmt"
	"os"

	"github.com/dikaio/scribes/pkg/cli"
)

func main() {
	// Initialize CLI
	app := cli.NewApp()

	// Run the app
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
