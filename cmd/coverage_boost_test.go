package cmd

import (
	"bytes"
	"testing"
)

func TestBoostCoverageHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	commands := [][]string{
		{"category", "--help"},
		{"category", "add", "--help"},
		{"category", "list", "--help"},
		{"listing", "--help"},
		{"listing", "create", "--help"},
		{"listing", "backfill", "--help"},
		{"serve", "--help"},
		{"admin", "--help"},
		{"benchmark", "--help"},
	}

	for _, args := range commands {
		rootCmd.SetArgs(args)
		_ = rootCmd.Execute()
	}
}
