package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Root command should not error when run with --help
	// We capture execution by args
	oldArgs := rootCmd.Args
	defer func() { rootCmd.Args = oldArgs }()

	// Just verify the root command exists
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "agbalumo", rootCmd.Use)

	// Verify subcommands
	commands := rootCmd.Commands()
	foundServe := false
	foundSeed := false
	for _, c := range commands {
		if c.Use == "serve" {
			foundServe = true
		}
		if c.Use == "seed" {
			foundSeed = true
		}
	}
	assert.True(t, foundServe, "serve command should be registered")
	assert.True(t, foundSeed, "seed command should be registered")
}
