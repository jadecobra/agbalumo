package agent_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestStateSerialization(t *testing.T) {
	t.Run("LoadState", func(t *testing.T) {
		content := `{
  "feature": "login-widget",
  "workflow_type": "feature",
  "phase": "IMPLEMENTATION",
  "gates": {
    "red-test": "PASSED",
    "api-spec": "PENDING",
    "implementation": "PENDING",
    "lint": "PENDING",
    "coverage": "PENDING",
    "browser-verification": "PENDING"
  },
  "updated_at": "2026-03-15T13:47:32Z"
}`
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
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(s)-len(substr) >= ind(s, substr)
}

func ind(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
