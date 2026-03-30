package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyCoverageOrchestration(t *testing.T) {
	// Since VerifyCoverage is highly coupled with hardcoded paths,
	// we use a temporary directory and change working directory.

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	// Case 1: Missing coverage profile
	if VerifyCoverage() {
		t.Error("VerifyCoverage should return false when profile is missing")
	}

	// Create necessary directories
	_ = os.MkdirAll(filepath.Join(".tester", "coverage"), 0755)
	_ = os.MkdirAll(".agents", 0755)

	// Case 2: Malformed profile
	covPath := filepath.Join(".tester", "coverage", "coverage.out")
	_ = os.WriteFile(covPath, []byte("invalid content"), 0644)
	if VerifyCoverage() {
		t.Error("VerifyCoverage should return false for malformed profile")
	}

	// Case 3: Success with thresholds file
	_ = os.WriteFile(covPath, []byte("mode: set\ngithub.com/jadecobra/agbalumo/main.go:1.0,2.0 1 1\n"), 0644)

	config := CoverageConfig{
		Thresholds: map[string]float64{"default": 50.0},
	}
	config.Signature = calculateCoverageSignature(&config)
	b, _ := json.Marshal(config)
	_ = os.WriteFile(filepath.Join(".agents", "coverage-thresholds.json"), b, 0644)

	if !VerifyCoverage() {
		t.Error("VerifyCoverage should return true when thresholds are met")
	}
}

func TestVerifyCoverage_LegacyThreshold(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	_ = os.MkdirAll(filepath.Join(".tester", "coverage"), 0755)
	_ = os.MkdirAll(".agents", 0755)

	covPath := filepath.Join(".tester", "coverage", "coverage.out")
	_ = os.WriteFile(covPath, []byte("mode: set\ngithub.com/jadecobra/agbalumo/main.go:1.0,2.0 1 1\n"), 0644)

	// Legacy threshold file
	_ = os.WriteFile(filepath.Join(".agents", "coverage-threshold"), []byte("95.5"), 0644)

	// 100% coverage should pass 95.5%
	if !VerifyCoverage() {
		t.Error("VerifyCoverage should pass with legacy threshold when met")
	}

	// Now fail legacy threshold
	_ = os.WriteFile(covPath, []byte("mode: set\ngithub.com/jadecobra/agbalumo/main.go:1.0,2.0 1 0\n"), 0644)
	if VerifyCoverage() {
		t.Error("VerifyCoverage should fail with legacy threshold when not met")
	}
}
