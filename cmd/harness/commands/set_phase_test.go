package commands

import (
	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"testing"
)

func TestSetPhaseCmd(t *testing.T) {
	cmd := SetPhaseCmd()

	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	state := &agent.State{
		Feature:      "test-feature",
		WorkflowType: "feature",
		Phase:        "IDLE",
	}
	_ = agent.SaveState(stateFile, state)
	defer func() { _ = util.SafeRemove(stateFile) }()

	// 1. Valid phases
	phases := []string{"RED", "GREEN", "REFACTOR", "IDLE"}
	for _, p := range phases {
		err := cmd.RunE(cmd, []string{p})
		if err != nil {
			t.Errorf("unexpected error for phase %s: %v", p, err)
		}
	}

	// 2. Invalid phase
	err := cmd.RunE(cmd, []string{"INVALID"})
	if err == nil {
		t.Errorf("expected error for invalid phase, got nil")
	}

	// 3. Command usage with --text
	flagText = true
	err = cmd.RunE(cmd, []string{"RED"})
	if err != nil {
		t.Fatalf("unexpected error for status with flagText: %v", err)
	}
}
