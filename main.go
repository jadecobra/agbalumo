package main

import (
	"os"

	"github.com/jadecobra/agbalumo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
