package maintenance

import (
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// Setup directory structure
	dirs := []string{
		"internal/repository",
		"internal/handler",
		"docs/adr",
		".agents/workflows",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(rootDir, d), 0750); err != nil {
			t.Fatal(err)
		}
	}

	// Create AGENTS.md files
	if err := os.WriteFile(filepath.Join(rootDir, "AGENTS.md"), []byte("Root AGENTS"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rootDir, "internal/repository/AGENTS.md"), []byte("Repo AGENTS"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create ADR
	if err := os.WriteFile(filepath.Join(rootDir, "docs/adr/2026-04-06-test.md"), []byte("# Test ADR\nThis mentions internal/repository"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create coding-standards.md
	codingContent := `## Strict Lessons
### CI & Infrastructure
* Infrastructure Lesson
### UI & Frontend
* Frontend Lesson
`
	if err := os.WriteFile(filepath.Join(rootDir, ".agents/workflows/coding-standards.md"), []byte(codingContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create invariants.json
	if err := os.WriteFile(filepath.Join(rootDir, ".agents/invariants.json"), []byte(`{"db_engine": "sqlite"}`), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		target         string
		expectedSubstr []string
	}{
		{
			name:   "repository domain",
			target: "internal/repository",
			expectedSubstr: []string{
				"📋 Session Context for: internal/repository",
				"📁 Local AGENTS.md (internal/repository/AGENTS.md)",
				"Repo AGENTS",
				"📁 Inherited AGENTS.md (AGENTS.md)",
				"Root AGENTS",
				"📚 Related ADRs:",
				"2026-04-06-test.md: Test ADR",
				"🔧 Invariants:",
				"db_engine: sqlite",
			},
		},
		{
			name:   "handler domain with lessons",
			target: "internal/handler",
			expectedSubstr: []string{
				"📋 Session Context for: internal/handler",
				"⚠️  Relevant Strict Lessons (UI & Frontend)",
				"Frontend Lesson",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := RunSessionContext(rootDir, filepath.Join(rootDir, tt.target))

			if errClose := w.Close(); errClose != nil {
				t.Fatal(errClose)
			}
			os.Stdout = old

			var buf strings.Builder
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			if err != nil {
				t.Errorf("RunSessionContext failed: %v", err)
			}

			for _, substr := range tt.expectedSubstr {
				if !strings.Contains(output, substr) {
					t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", substr, output)
				}
			}
		})
	}
}
