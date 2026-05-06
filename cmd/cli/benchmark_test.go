package cli

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

	warmup = true
	BenchmarkCmd.Run(BenchmarkCmd, []string{dbPath})

	// Make sure warmup flag successfully resets
	warmup = false // Reset manually for testing environment pollution
}
