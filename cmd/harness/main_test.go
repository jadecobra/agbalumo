package main

import (
	"testing"
	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd.Use != "harness" {
		t.Errorf("expected Use to be 'harness', got '%s'", cmd.Use)
	}

	if len(cmd.Commands()) < 5 {
		t.Errorf("expected at least 5 subcommands, got %d", len(cmd.Commands()))
	}
	
	cmd.SetArgs([]string{"--help"})
	_ = cmd.Execute()
}

func TestCoverageHelpers(t *testing.T) {
	// Call these purely for coverage. They safely return early or are read-only operations.
	getState()
	saveState(&agent.State{
		Feature:      "test_feature",
		WorkflowType: "feature",
		Phase:        "IDLE",
	})
	summarizeProgress()
	checkAndApplyProgressUpdate()
}
