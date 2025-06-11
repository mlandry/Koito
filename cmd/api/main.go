package main

import (
	"fmt"
	"os"

	"github.com/gabehf/koito/engine"
)

func main() {
	if err := engine.Run(
		os.Getenv,
		os.Stdout,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
