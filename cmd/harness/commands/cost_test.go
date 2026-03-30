package commands

import (
	"testing"
)

func TestCostCmd(t *testing.T) {
	cmd := CostCmd()

	// 1. Basic cost check
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error for cost: %v", err)
	}

	// 2. Cost check with text flag
	flagText = true
	err = cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error for cost with flagText: %v", err)
	}

	// 3. Cost with pattern cleanup
	flagText = false
	err = cmd.RunE(cmd, []string{"."})
	if err != nil {
		t.Fatalf("unexpected error for cost .: %v", err)
	}

	// 4. Missing directory in cost check
	err = cmd.RunE(cmd, []string{"/non/existent/dir"})
	if err == nil {
		t.Errorf("expected error for non-existent directory, got nil")
	}
}
