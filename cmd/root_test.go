package cmd

import (
	"testing"
)

func TestExecute(t *testing.T) {
	// Test that Execute returns nil (success) when run without args (or default args)
	// Since rootCmd.Run is empty/default, it should just succeed.
	// We are testing the wiring.
	if err := Execute(); err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
}
