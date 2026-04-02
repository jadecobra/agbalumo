package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestStateLoad(t *testing.T) {
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
    "template-drift": "PENDING",
    "security-static": "PENDING"
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
}
