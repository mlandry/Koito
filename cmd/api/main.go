package main

import (
	"fmt"
	"os"

	"github.com/gabehf/koito/engine"
)

var Version = "dev"

func main() {
	if err := engine.Run(
		os.Getenv,
		os.Stdout,
		Version,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
