package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestStressCmd_Success(t *testing.T) {
	// Create a temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "stress_test.db")

	// Ensure DB is removed even on failure
	defer func() {
		_ = os.Remove(dbPath)
	}()

	stressCount = 50
	StressCmd.Run(StressCmd, []string{dbPath})

	// Verify database state
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	listings, _, err := repo.FindAll(context.Background(), "", "", "", 0.0, 0.0, 0.0, "", "", false, 100, 0)
	if err != nil {
		t.Fatalf("failed to query test database: %v", err)
	}

	if len(listings) != 50 {
		t.Errorf("expected 50 listings inserted, got %d", len(listings))
	}
}
