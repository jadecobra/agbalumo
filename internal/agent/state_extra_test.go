package agent

import (
	"os"
	"testing"
)

func TestStateExtraCoverage(t *testing.T) {
	// 1. Test IsNotExist
	if !IsNotExist(os.ErrNotExist) {
		t.Error("expected IsNotExist to be true")
	}

	// 2. Test SaveState and calculateSignature
	testPath := "test_state.json"
	defer func() {
		_ = os.Remove(testPath)
	}()

	state := &State{
		Feature: "test-feature-x",
	}

	err := SaveState(testPath, state)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Test LoadState
	loaded, errLoad := LoadState(testPath)
	if errLoad != nil {
		t.Fatal(errLoad)
	}
	if loaded.Feature != "test-feature-x" {
		t.Error("feature mismatch")
	}

	// 4. Test Error paths
	t.Run("LoadState_NotFound", func(t *testing.T) {
		_, errNotFound := LoadState("does-not-exist.json")
		if errNotFound == nil {
			t.Error("expected error loading non-existent file")
		}
	})

	t.Run("SaveState_WriteError", func(t *testing.T) {
		// Try to save to a directory that is actually a file
		err := SaveState(testPath+"/invalid/path", state)
		if err == nil {
			t.Error("expected error saving to invalid path, got none")
		}
	})

	t.Run("CalculateSignature_Consistency", func(t *testing.T) {
		s1 := &State{Feature: "feat1", WorkflowType: "refactor"}
		sig1 := calculateSignature(s1)

		s2 := &State{Feature: "feat1", WorkflowType: "refactor"}
		sig2 := calculateSignature(s2)

		if sig1 != sig2 {
			t.Errorf("expected consistent signatures for same state, got %s and %s", sig1, sig2)
		}

		s3 := &State{Feature: "feat1", WorkflowType: "refactor", Signature: "something"}
		sig3 := calculateSignature(s3)
		if sig1 != sig3 {
			t.Error("calculateSignature should ignore existing Signature field")
		}

		// Known hash check (optional but good for regression)
		// Since UpdatedAt is not set yet in s1, it's predictable.
		// Feature: feat1, WorkflowType: refactor, others empty.
	})
}
