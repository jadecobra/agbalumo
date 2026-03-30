package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestCommandErrors(t *testing.T) {
	// Clean up state file before and after
	_ = os.Remove(".agents/state.json")
	defer func() { _ = os.Remove(".agents/state.json") }()

	t.Run("GateCmd Invalid Status", func(t *testing.T) {
		cmd := GateCmd()
		_, err := executeCommand(cmd, "lint", "INVALID")
		if err == nil {
			t.Error("expected error for invalid status, got nil")
		}
	})

	t.Run("GateCmd Manual Bypass", func(t *testing.T) {
		cmd := GateCmd()
		_, err := executeCommand(cmd, "lint", "PASS")
		if err == nil {
			t.Errorf("expected error for manual bypass of lint gate, got nil")
		}
	})

	t.Run("SetPhaseCmd Invalid Phase", func(t *testing.T) {
		cmd := SetPhaseCmd()
		_, err := executeCommand(cmd, "INVALID")
		if err == nil {
			t.Error("expected error for invalid phase, got nil")
		}
	})

	t.Run("StatusCmd Success", func(t *testing.T) {
		// Create a dummy state file
		root := NewRootCmd()
		_, _ = executeCommand(root, "init", "test-feat", "feature")

		cmd := StatusCmd()
		_, err := executeCommand(cmd)
		if err != nil {
			t.Fatalf("unexpected error for status command: %v", err)
		}
	})
}
