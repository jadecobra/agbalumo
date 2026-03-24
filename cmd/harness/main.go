package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/cmd/harness/commands"
)

func main() {
	cmd := commands.NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
