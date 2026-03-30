package commands

import (
	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"testing"
)

func TestGateCmd(t *testing.T) {
	cmd := GateCmd()

	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	state := &agent.State{
		Feature: "test-feature",
		Phase:   "RED",
	}
	_ = agent.SaveState(stateFile, state)
	defer func() { _ = util.SafeRemove(stateFile) }()

	// 1. Success case: PENDING
	err := cmd.RunE(cmd, []string{"browser-verification", "PENDING"})
	if err != nil {
		t.Fatalf("unexpected error for PENDING: %v", err)
	}

	// 2. Success case: PASS (allowed for browser-verification)
	err = cmd.RunE(cmd, []string{"browser-verification", "PASS"})
	if err != nil {
		t.Fatalf("unexpected error for PASS: %v", err)
	}

	// 2b. TEST FAIL status for browser-verification
	flagText = false
	err = cmd.RunE(cmd, []string{"browser-verification", "FAIL"})
	if err != nil {
		t.Fatalf("unexpected error for FAIL: %v", err)
	}

	// 3. Failure case: manual bypass not allowed
	failures := []string{"red-test", "api-spec", "implementation", "lint", "coverage"}
	for _, f := range failures {
		err = cmd.RunE(cmd, []string{f, "PASS"})
		if err == nil {
			t.Errorf("expected error for manual bypass of %s, got nil", f)
		}
	}

	// 4. Invalid status
	err = cmd.RunE(cmd, []string{"browser-verification", "INVALID"})
	if err == nil {
		t.Errorf("expected error for invalid status, got nil")
	}

	// 5. Unknown gate
	err = cmd.RunE(cmd, []string{"unknown-gate", "PASS"})
	if err == nil {
		t.Errorf("expected error for unknown gate, got nil")
	}
}
