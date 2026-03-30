package main

import (
	"os"

	"github.com/chuck/openspec-go/internal/cmd"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
