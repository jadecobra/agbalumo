package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunPerformanceAudit(t *testing.T) {
	// Create dummy assets for testing
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "ui/static/css")
	jsDir := filepath.Join(tmpDir, "ui/static/js")
	repoDir := filepath.Join(tmpDir, "internal/repository/sqlite")

	os.MkdirAll(cssDir, 0755)
	os.MkdirAll(jsDir, 0755)
	os.MkdirAll(repoDir, 0755)

	os.WriteFile(filepath.Join(cssDir, "output.css"), make([]byte, 100*1024), 0644) // 100KB < 150KB
	os.WriteFile(filepath.Join(jsDir, "app.js"), make([]byte, 10*1024), 0644)       // 10KB < 50KB

	sqliteContent := `
		package sqlite
		func setup() {
			db.Exec("PRAGMA journal_mode=WAL;")
			db.Exec("PRAGMA busy_timeout=5000;")
			db.SetMaxOpenConns(100)
		}
	`
	os.WriteFile(filepath.Join(repoDir, "sqlite.go"), []byte(sqliteContent), 0644)

	err := RunPerformanceAudit(tmpDir)
	if err != nil {
		t.Fatalf("Performance audit failed: %v", err)
	}
}
