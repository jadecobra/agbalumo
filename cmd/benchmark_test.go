package cmd

import (
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestBenchmarkCmd_Success(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark_test.db")

	// Create and initialize DB schema
	_, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	rootCmd.SetArgs([]string{"benchmark", dbPath, "--warmup"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing benchmark command: %v", err)
	}

	// Make sure warmup flag successfully resets
	warmup = false // Reset manually for testing environment pollution
}
