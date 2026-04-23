package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSessionContext(t *testing.T) {
	rootDir, err := os.MkdirTemp("", "session-context-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(rootDir)
	}()

	// Fix shadowing by using err = instead of err :=
	var errFS error

	errFS = os.MkdirAll(filepath.Join(rootDir, "internal/repository"), 0750)
	if errFS != nil {
		t.Fatal(errFS)
	}
	errFS = os.MkdirAll(filepath.Join(rootDir, "docs/adr"), 0750)
	if errFS != nil {
		t.Fatal(errFS)
	}
	errFS = os.MkdirAll(filepath.Join(rootDir, ".agents/workflows"), 0750)
	if errFS != nil {
		t.Fatal(errFS)
	}

	// Create AGENTS.md
	errFS = os.WriteFile(filepath.Join(rootDir, "AGENTS.md"), []byte("Root AGENTS"), 0600)
	if errFS != nil {
		t.Fatal(errFS)
	}
	errFS = os.WriteFile(filepath.Join(rootDir, "internal/repository/AGENTS.md"), []byte("Repo AGENTS"), 0600)
	if errFS != nil {
		t.Fatal(errFS)
	}

	// Create ADR
	errFS = os.WriteFile(filepath.Join(rootDir, "docs/adr/2026-04-06-test.md"), []byte("# Test ADR\nThis mentions internal/repository"), 0600)
	if errFS != nil {
		t.Fatal(errFS)
	}

	// Create coding-standards.md
	errFS = os.WriteFile(filepath.Join(rootDir, ".agents/workflows/coding-standards.md"), []byte("## Strict Lessons\n### CI & Infrastructure\n* Lesson 1"), 0600)
	if errFS != nil {
		t.Fatal(errFS)
	}

	// Create invariants.json
	errFS = os.WriteFile(filepath.Join(rootDir, ".agents/invariants.json"), []byte(`{"db_engine": "sqlite"}`), 0600)
	if errFS != nil {
		t.Fatal(errFS)
	}

	// Run it
	err = RunSessionContext(rootDir, filepath.Join(rootDir, "internal/repository"))
	if err != nil {
		t.Errorf("RunSessionContext failed: %v", err)
	}
}
