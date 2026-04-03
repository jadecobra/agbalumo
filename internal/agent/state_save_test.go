package agent_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestStateSave(t *testing.T) {
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
				SecurityStatic:      agent.GatePending,
				VibeCheck:           agent.GatePending,
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

		loadedState, err := agent.LoadState(statePath)
		if err != nil {
			t.Fatalf("expected no error loading written state, got %v", err)
		}
		if loadedState.UpdatedAt.Equal(time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)) {
			if time.Since(loadedState.UpdatedAt) > time.Minute {
				t.Errorf("expected UpdatedAt to be updated to now, got %v", loadedState.UpdatedAt)
			}
		}
	})

	t.Run("SaveState_Permissions", func(t *testing.T) {
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

		if info.Mode().Perm() != 0644 {
			t.Errorf("expected perm 0644, got %o", info.Mode().Perm())
		}
	})
}
