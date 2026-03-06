package main

import (
	"os"

	"github.com/pavelvarganov/lcli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
