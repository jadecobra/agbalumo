package agent_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestStateSecurity(t *testing.T) {
	t.Run("LoadState_InvalidCasingSignature", func(t *testing.T) {
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
    "template-drift": "PENDING",
    "security-static": "PENDING",
    "vibe-check": "PENDING"
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

	t.Run("LoadState_StructuralMismatch_ExtraField", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.json")
		state := &agent.State{Feature: "structural"}
		err := agent.SaveState(statePath, state)
		if err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}

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
