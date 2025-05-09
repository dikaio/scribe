package main

import (
	"fmt"
	"os"

	"github.com/dikaio/scribe/internal/release"
)

func main() {
	// Pass all arguments to the release package
	if err := release.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}