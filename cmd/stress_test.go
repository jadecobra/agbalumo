package cmd

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

	// Since rootCmd executes stress via `rootCmd.Execute()`,
	// we set arguments directly to sub-command.
	// `ResolveSeedConfig` logic implies if `args` has an element, it might use it as db path.
	// `agbalumo stress custom.db --count 50`
	rootCmd.SetArgs([]string{"stress", dbPath, "--count", "50"})

	// Execute command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing stress command: %v", err)
	}

	// Verify database state
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	listings, _, err := repo.FindAll(context.Background(), "", "", "", "", false, 100, 0)
	if err != nil {
		t.Fatalf("failed to query test database: %v", err)
	}

	if len(listings) != 50 {
		t.Errorf("expected 50 listings inserted, got %d", len(listings))
	}
}
