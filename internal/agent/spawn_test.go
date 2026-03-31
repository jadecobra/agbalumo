package agent

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSpawnAgent_Integration(t *testing.T) {
	// Setup a temporary directory for out-of-process verification
	tmpDir, err := os.MkdirTemp("", "spawn-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock antigravity script that writes to a file
	mockScript := filepath.Join(tmpDir, "antigravity")
	content := `#!/bin/bash
sleep 1
touch ` + tmpDir + `/worked
`
	if err := os.WriteFile(mockScript, []byte(content), 0755); err != nil {
		t.Fatal(err)
	}

	// Update PATH so our mock script is found
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tmpDir+":"+oldPath)

	start := time.Now()
	err = SpawnAgent("some prompt")
	duration := time.Since(start)

	if err != nil {
		t.Errorf("SpawnAgent failed: %v", err)
	}

	// It should be extremely fast because it's detached
	if duration > 500*time.Millisecond {
		t.Errorf("SpawnAgent blocked for %v, should be asynchronous", duration)
	}

	// Wait a bit to see if the background process actually worked
	time.Sleep(2 * time.Second)
	if _, err := os.Stat(filepath.Join(tmpDir, "worked")); os.IsNotExist(err) {
		t.Errorf("Background process did not seem to run (worked file missing)")
	}
}
