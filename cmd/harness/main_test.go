package main

import (
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd.Use != "harness" {
		t.Errorf("expected Use to be 'harness', got '%s'", cmd.Use)
	}

	if len(cmd.Commands()) < 5 {
		t.Errorf("expected at least 5 subcommands, got %d", len(cmd.Commands()))
	}
}
