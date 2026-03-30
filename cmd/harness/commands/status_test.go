package commands

import (
	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"testing"
)

func TestStatusCmd(t *testing.T) {
	cmd := StatusCmd()

	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	state := &agent.State{
		Feature:      "test-feature",
		WorkflowType: "feature",
		Phase:        "REFACTOR",
		Gates: agent.Gates{
			RedTest: "PASSED",
		},
	}
	_ = agent.SaveState(stateFile, state)
	defer func() { _ = util.SafeRemove(stateFile) }()

	// 1. Basic status check
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error for status: %v", err)
	}

	// 2. Status check with text flag
	flagText = true
	err = cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error for status with flagText: %v", err)
	}
}
