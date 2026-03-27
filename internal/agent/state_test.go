package agent_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestStateSerialization(t *testing.T) {
	t.Run("LoadState", func(t *testing.T) {
		content := `{
  "_DO_NOT_EDIT_": "MANUAL EDITS WILL INVALIDATE SIGNATURE. USE ./scripts/agent-exec.sh TO MANAGE STATE.",
  "feature": "login-widget",
  "workflow_type": "feature",
  "phase": "IMPLEMENTATION",
  "gates": {
    "red-test": "PASSED",
    "api-spec": "PENDING",
    "implementation": "PENDING",
    "lint": "PENDING",
    "coverage": "PENDING",
    "browser-verification": "PENDING",
    "template-drift": "PENDING"
  },
  "updated_at": "2026-03-15T13:47:32Z",
  "signature": ""
}
`
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		err := os.WriteFile(statePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		state, err := agent.LoadState(statePath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if state.Feature != "login-widget" {
			t.Errorf("expected feature login-widget, got %s", state.Feature)
		}
		if state.Phase != "IMPLEMENTATION" {
			t.Errorf("expected phase IMPLEMENTATION, got %s", state.Phase)
		}
		if state.Gates.RedTest != agent.GatePassed {
			t.Errorf("expected red-test PASSED, got %s", state.Gates.RedTest)
		}
		if state.Gates.ApiSpec != agent.GatePending {
			t.Errorf("expected api-spec PENDING, got %s", state.Gates.ApiSpec)
		}
	})

	t.Run("SaveState", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")

		state := &agent.State{
			Feature:      "new-feature",
			WorkflowType: "feature",
			Phase:        "IDLE",
			Gates: agent.Gates{
				RedTest:             agent.GatePassed,
				ApiSpec:             agent.GatePassed,
				Implementation:      agent.GatePending,
				Lint:                agent.GatePending,
				Coverage:            agent.GatePending,
				BrowserVerification: agent.GatePending,
				TemplateDrift:       agent.GatePending,
			},
			UpdatedAt: time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC),
		}

		err := agent.SaveState(statePath, state)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		b, err := os.ReadFile(statePath)
		if err != nil {
			t.Fatalf("failed to read written file: %v", err)
		}

		writtenContent := string(b)
		if !contains(writtenContent, `"feature": "new-feature"`) {
			t.Errorf("expected written content to contain feature name, got: %s", writtenContent)
		}
		if !contains(writtenContent, `"red-test": "PASSED"`) {
			t.Errorf("expected written content to contain red-test passed, got: %s", writtenContent)
		}
		// ensure updated_at was modified by SaveState to be recent
		loadedState, err := agent.LoadState(statePath)
		if err != nil {
			t.Fatalf("expected no error loading written state, got %v", err)
		}
		// SaveState should update the UpdatedAt timestamp to a newer time.
		if loadedState.UpdatedAt.Equal(time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)) {
			// Actually depending on implementation it might or might not update it.
			// The protocol states we just want robust read/write. I'll require SaveState to update UpdatedAt.
			if time.Since(loadedState.UpdatedAt) > time.Minute {
				t.Errorf("expected UpdatedAt to be updated to now, got %v", loadedState.UpdatedAt)
			}
		}
	})
	
	t.Run("LoadState_NotFound", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "nonexistent.json")
		
		_, err := agent.LoadState(statePath)
		if err == nil {
			t.Fatalf("expected error loading non-existent file")
		}
		if !agent.IsNotExist(err) && !os.IsNotExist(err) {
			t.Errorf("expected IsNotExist error, got %v", err)
		}
	})

	t.Run("LoadState_InvalidJSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(statePath, []byte(`{invalid json`), 0644)
		if err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		_, err = agent.LoadState(statePath)
		if err == nil {
			t.Fatal("expected error for invalid JSON, got none")
		}
	})

	t.Run("LoadState_InvalidJSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(statePath, []byte(`{invalid json`), 0644)
		if err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		_, err = agent.LoadState(statePath)
		if err == nil {
			t.Fatal("expected error for invalid JSON, got none")
		}
	})

	t.Run("LoadState_InvalidCasingSignature", func(t *testing.T) {
		// Valid signature for lower-case "feature"
		content := `{
  "_DO_NOT_EDIT_": "MANUAL EDITS WILL INVALIDATE SIGNATURE. USE ./scripts/agent-exec.sh TO MANAGE STATE.",
  "Feature": "login-widget",
  "workflow_type": "feature",
  "phase": "IMPLEMENTATION",
  "gates": {
    "red-test": "PASSED",
    "api-spec": "PENDING",
    "implementation": "PENDING",
    "lint": "PENDING",
    "coverage": "PENDING",
    "browser-verification": "PENDING",
    "template-drift": "PENDING"
  },
  "updated_at": "2026-03-15T13:47:32Z",
  "signature": "ce73463aee24cd8c2ba43dc759bc65bcba172b8c9d0b674bfa0e6f3b55c6ce8e"
}`
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		err := os.WriteFile(statePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		_, err = agent.LoadState(statePath)
		if err == nil {
			t.Fatalf("expected error due to casing exploit, got none")
		}
		if !contains(err.Error(), "structural mismatch") {
			t.Errorf("expected structural mismatch error, got: %v", err)
		}
	})

	t.Run("LoadState_InvalidSignature", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		state := &agent.State{Feature: "original"}
		err := agent.SaveState(statePath, state)
		if err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}

		// Read and tamper with value but keep signature
		b, _ := os.ReadFile(statePath)
		tampered := string(b)
		tampered = replaceValue(tampered, `"feature": "original"`, `"feature": "tampered"`)
		
		err = os.WriteFile(statePath, []byte(tampered), 0644)
		if err != nil {
			t.Fatalf("failed to write tampered file: %v", err)
		}

		_, err = agent.LoadState(statePath)
		if err == nil {
			t.Fatal("expected error for tampered signature, got none")
		}
		if !contains(err.Error(), "ANTI-CHEAT TRIGGERED") || contains(err.Error(), "structural mismatch") {
			t.Errorf("expected signature anti-cheat error, got: %v", err)
		}
	})

	t.Run("LoadState_Permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		state := &agent.State{Feature: "perm-test"}

		err := agent.SaveState(statePath, state)
		if err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}

		info, err := os.Stat(statePath)
		if err != nil {
			t.Fatalf("Stat failed: %v", err)
		}

		if info.Mode().Perm() != 0600 {
			t.Errorf("expected perm 0600, got %o", info.Mode().Perm())
		}
	})

	t.Run("LoadState_ValidWithSignature", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		state := &agent.State{
			Feature:      "valid-feature",
			WorkflowType: "feature",
			Phase:        "READY",
		}
		
		err := agent.SaveState(statePath, state)
		if err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}

		loaded, err := agent.LoadState(statePath)
		if err != nil {
			t.Fatalf("LoadState failed: %v", err)
		}
		if loaded.Feature != "valid-feature" {
			t.Errorf("expected feature valid-feature, got %s", loaded.Feature)
		}
		if loaded.Signature == "" {
			t.Error("expected signature to be set by SaveState")
		}
	})

	t.Run("LoadState_StructuralMismatch_ExtraField", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		state := &agent.State{Feature: "structural"}
		err := agent.SaveState(statePath, state)
		if err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}

		// Manually add a field to the JSON
		b, _ := os.ReadFile(statePath)
		var m map[string]interface{}
		if uerr := json.Unmarshal(b, &m); uerr != nil {
			t.Fatalf("Unmarshal failed: %v", uerr)
		}
		m["extra"] = "field"
		b2, _ := json.MarshalIndent(m, "", "  ")
		b2 = append(b2, '\n')
		if werr := os.WriteFile(statePath, b2, 0644); werr != nil {
			t.Fatalf("WriteFile failed: %v", werr)
		}

		_, err = agent.LoadState(statePath)
		if err == nil {
			t.Fatal("expected error for structural mismatch (extra field), got none")
		}
		if !contains(err.Error(), "structural mismatch") {
			t.Errorf("expected structural mismatch error, got: %v", err)
		}
	})
}

func replaceValue(s, old, new string) string {
	res := ""
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			res = s[:i] + new + s[i+len(old):]
			break
		}
	}
	return res
}

func contains(s, substr string) bool {
	return ind(s, substr) != -1
}

func ind(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
