package commands

import (
	"github.com/jadecobra/agbalumo/internal/util"
	"testing"
)

func TestRootFunctions(t *testing.T) {
	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	// 1. Test getState (missing file)
	_ = util.SafeRemove(stateFile)
	state, err := getState()
	if err != nil {
		t.Fatalf("unexpected error for missing state: %v", err)
	}
	if state.Phase != "" && state.Phase != "IDLE" {
		t.Errorf("expected empty or IDLE phase for new state, got %s", state.Phase)
	}

	// 2. Test saveState
	state.Feature = "test-save"
	err = saveState(state)
	if err != nil {
		t.Fatalf("unexpected error for saveState: %v", err)
	}

	// 3. Test getState (existing file)
	state2, err := getState()
	if err != nil {
		t.Fatalf("unexpected error for reading back state: %v", err)
	}
	if state2.Feature != "test-save" {
		t.Errorf("expected feature test-save, got %s", state2.Feature)
	}

	// 4. Test summarizeProgress (missing file)
	pfile := ".tester/tasks/progress.md"
	pbak := pfile + ".bak"
	if _, pErr := util.SafeStat(pfile); pErr == nil {
		_ = util.SafeRename(pfile, pbak)
		defer func() { _ = util.SafeRename(pbak, pfile) }()
	}
	_ = util.SafeRemove(pfile)
	err = summarizeProgress()
	if err == nil {
		t.Error("expected error for missing progress file, got nil")
	}

	// 5. Test summarizeProgress (valid file)
	_ = util.SafeWriteFile(pfile, []byte("# Category\n- [ ] Task 1\n- [x] Task 2\n"))
	defer func() { _ = util.SafeRemove(pfile) }()
	flagText = true
	err = summarizeProgress()
	if err != nil {
		t.Errorf("unexpected error for valid progress file: %v", err)
	}

	// 6. Test checkAndApplyProgressUpdate (missing file)
	updateFile := ".tester/tasks/pending_update.md"
	_ = util.SafeRemove(updateFile)
	err = checkAndApplyProgressUpdate()
	if err != nil {
		t.Fatalf("unexpected error for missing progress update: %v", err)
	}

	// 7. Test checkAndApplyProgressUpdate (valid file)
	_ = util.SafeWriteFile(updateFile, []byte("# New Category\n- [x] Done Task\n"))
	defer func() { _ = util.SafeRemove(updateFile) }()

	err = checkAndApplyProgressUpdate()
	if err != nil {
		t.Fatalf("unexpected error for progress update: %v", err)
	}

	// 8. Test summarizeProgress (non-text mode)
	flagText = false
	err = summarizeProgress()
	if err != nil {
		t.Fatalf("unexpected error for non-text summary: %v", err)
	}
}

func TestPrintJSON(t *testing.T) {
	// printJSON uses os.Stdout, we just verify it doesn't panic
	printJSON(true, "test", map[string]any{"key": "val"}, nil)
	printJSON(false, "test", nil, []string{"warn1"})
}

func TestGetState_InvalidJSON(t *testing.T) {
	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	_ = util.SafeWriteFile(stateFile, []byte("invalid json"))
	defer func() { _ = util.SafeRemove(stateFile) }()

	_, err := getState()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
