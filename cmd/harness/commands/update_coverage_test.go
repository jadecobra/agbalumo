package commands

import (
	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"testing"
)

func TestUpdateCoverageCmd(t *testing.T) {
	cmd := UpdateCoverageCmd()

	thresholdsFile := ".agents/coverage-thresholds.json"
	backupFile := thresholdsFile + ".bak"
	if _, err := util.SafeStat(thresholdsFile); err == nil {
		_ = util.SafeRename(thresholdsFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, thresholdsFile) }()
	}
	_ = agent.SaveThresholds(thresholdsFile, map[string]float64{})
	defer func() { _ = util.SafeRemove(thresholdsFile) }()

	// 1. Initial creation
	cmd.SetArgs([]string{"all", "0.5"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error for initial update: %v", err)
	}

	// 2. Increase threshold
	cmd.SetArgs([]string{"all", "0.6"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error for increase: %v", err)
	}

	// 3. Lowering threshold protection
	cmd.SetArgs([]string{"all", "0.4"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("expected error for lowering threshold, got nil")
	}

	// 4. Invalid threshold
	cmd.SetArgs([]string{"all", "invalid"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("expected error for invalid threshold, got nil")
	}

	// 5. Missing arguments (vaguely testing cobra validation)
	cmd.SetArgs([]string{"all"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("expected error for missing args, got nil")
	}
}
