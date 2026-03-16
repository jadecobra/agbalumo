package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	cmd := NewRootCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "harness", cmd.Use)
}

func TestSubCommandsExist(t *testing.T) {
	cmd := NewRootCmd()

	expectedCommands := []string{"init", "set-phase", "status", "gate"}

	// Extract the names of registered subcommands
	var actualCommands []string
	for _, subcmd := range cmd.Commands() {
		actualCommands = append(actualCommands, subcmd.Name())
	}

	for _, expected := range expectedCommands {
		assert.Contains(t, actualCommands, expected, "Root command should contain the '%s' subcommand", expected)
	}
}

func TestCommandExecution_Help(t *testing.T) {
    cmd := NewRootCmd()
    b := bytes.NewBufferString("")
    cmd.SetOut(b)
    cmd.SetArgs([]string{"--help"})
    
    err := cmd.Execute()
    assert.NoError(t, err)
    
    out := b.String()
    assert.Contains(t, out, "harness")
    assert.Contains(t, out, "init")
    assert.Contains(t, out, "set-phase")
    assert.Contains(t, out, "status")
    assert.Contains(t, out, "gate")
}
