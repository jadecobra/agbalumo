package seeder_test

import (
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

// setupSeeder initialized a test repo and a temp config path.
func setupSeeder(t *testing.T) (*sqlite.SQLiteRepository, string, func()) {
	t.Helper()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "categories.json")
	return repo, configPath, func() { _ = repo.Close() }
}
