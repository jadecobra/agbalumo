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
	_, errNotFound := LoadState("does-not-exist.json")
	if errNotFound == nil {
		t.Error("expected error loading non-existent file")
	}
}
