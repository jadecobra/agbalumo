package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestUpdateCoverageCmd_Basic(t *testing.T) {
	cmd := &cobra.Command{Use: "update-coverage"}
	if cmd.Use != "update-coverage" {
		t.Errorf("expected update-coverage, got %s", cmd.Use)
	}
}

func TestExecuteUpdateCoverage(t *testing.T) {
	root := &cobra.Command{Use: "harness"}
	update := UpdateCoverageCmd()
	root.AddCommand(update)

	b := bytes.NewBufferString("")
	root.SetOut(b)
	root.SetArgs([]string{"update-coverage", "--help"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}
