package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
)

func TestChaosCommand(t *testing.T) {
	// Setup temporary state file for isolation
	tmpDir, err := os.MkdirTemp("", "chaos-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origStateFile := StateFile
	StateFile = filepath.Join(tmpDir, "state.json")
	defer func() { StateFile = origStateFile }()

	state := &agent.State{
		Feature: "chaos-test",
		Phase:   "RED",
	}
	if sErr := saveState(state); sErr != nil {
		t.Fatalf("failed to setup state: %v", sErr)
	}

	// 1. Test --state-corrupt
	cmd := ChaosCmd()
	cmd.SetArgs([]string{"--state-corrupt"})
	if cErr := cmd.Execute(); cErr != nil {
		t.Fatalf("chaos --state-corrupt failed: %v", cErr)
	}

	// Verify state file was written with corrupted signature
	data, err := util.SafeReadFile(StateFile)
	if err != nil {
		t.Fatalf("failed to read state: %v", err)
	}
	var loaded agent.State
	if uErr := json.Unmarshal(data, &loaded); uErr != nil {
		t.Fatalf("failed to unmarshal state: %v", uErr)
	}
	if !strings.HasSuffix(loaded.Signature, "_CORRUPTED") {
		t.Errorf("expected signature to end with _CORRUPTED, got %s", loaded.Signature)
	}

	// Verify LoadState fails as expected
	if _, lErr := agent.LoadState(StateFile); lErr == nil {
		t.Error("expected LoadState to fail after corruption, but it succeeded")
	} else if !strings.Contains(lErr.Error(), "ANTI-CHEAT TRIGGERED") {
		t.Errorf("expected anti-cheat error, got: %v", lErr)
	}

	// 2. Test --env-wipe
	envTmpDir := filepath.Join(tmpDir, "env_tmp")
	_ = util.SafeMkdir(envTmpDir)
	testFile := filepath.Join(envTmpDir, "chaos_wipe_test.txt")
	_ = util.SafeWriteFile(testFile, []byte("wipe me"))

	// Since --env-wipe is hardcoded to .tester/tmp, we'll verify it via a custom path if we can,
	// but the implementation uses a literal ".tester/tmp".
	// For testing, we'll just verify the command executes without error when the dir doesn't exist or is empty.
	cmd = ChaosCmd()
	cmd.SetArgs([]string{"--env-wipe"})
	if eErr := cmd.Execute(); eErr != nil {
		t.Fatalf("chaos --env-wipe failed: %v", eErr)
	}

	// 3. Test --test-sabotage
	targetTestFile := filepath.Join(tmpDir, "sabotage_test.go")
	testContent := "package commands\nimport \"testing\"\nfunc TestDummy(t *testing.T) {\n\t// dummy test\n}\n"
	_ = util.SafeWriteFile(targetTestFile, []byte(testContent))

	// Temporarily change directory to tmpDir so filepath.Glob finds our file
	origDir, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origDir) }()

	cmd = ChaosCmd()
	cmd.SetArgs([]string{"--test-sabotage"})
	if sErr := cmd.Execute(); sErr != nil {
		t.Fatalf("chaos --test-sabotage failed: %v", sErr)
	}

	content, err := util.SafeReadFile("sabotage_test.go")
	if err != nil {
		t.Fatalf("failed to read sabotaged file: %v", err)
	}

	if !strings.Contains(string(content), "CHAOS_SABOTAGE") {
		t.Error("expected sabotaged file to contain CHAOS_SABOTAGE, but it doesn't")
	}
}
