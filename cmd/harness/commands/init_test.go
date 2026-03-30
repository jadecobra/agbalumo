package commands

import (
	"github.com/jadecobra/agbalumo/internal/util"
	"testing"
)

func TestInitCmd(t *testing.T) {
	cmd := InitCmd()

	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	// 1. Success case: feature only
	cmd.SetArgs([]string{"test-feature"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error for init feature: %v", err)
	}

	// 2. Success case: feature + workflow
	cmd.SetArgs([]string{"test-bugfix", "bugfix"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error for init bugfix: %v", err)
	}

	// 3. Success case: feature + refactor + --text
	flagText = true
	cmd.SetArgs([]string{"test-refactor", "refactor"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error for init refactor: %v", err)
	}
	flagText = false

	// 4. Invalid workflow type
	cmd.SetArgs([]string{"test-feature", "invalid"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("expected error for invalid workflow type, got nil")
	}

	// 5. Missing feature
	cmd.SetArgs([]string{})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("expected error for missing args, got nil")
	}
}
