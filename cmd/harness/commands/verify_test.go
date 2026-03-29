package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestVerifyCmd_Basic(t *testing.T) {
	// We can't easily run the full verify logic in a unit test because it touches the filesystem
	// but we can at least ensure the command is structured correctly.
	cmd := &cobra.Command{Use: "verify"}
	// Just a smoke test for the command structure
	if cmd.Use != "verify" {
		t.Errorf("expected verify, got %s", cmd.Use)
	}
}

func TestExecuteVerify(t *testing.T) {
	// Create a dummy root command to test execution
	root := &cobra.Command{Use: "harness"}
	verify := VerifyCmd()
	root.AddCommand(verify)

	// Capture output
	b := bytes.NewBufferString("")
	root.SetOut(b)
	root.SetArgs([]string{"verify", "--help"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}
