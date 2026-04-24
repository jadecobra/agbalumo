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

	setupSessionContextTestEnv(t, rootDir)

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
			output := captureStdout(t, func() error {
				return RunSessionContext(rootDir, filepath.Join(rootDir, tt.target))
			})

			for _, substr := range tt.expectedSubstr {
				if !strings.Contains(output, substr) {
					t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", substr, output)
				}
			}
		})
	}
}

func setupSessionContextTestEnv(t *testing.T, rootDir string) {
	t.Helper()
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
	writeFile(t, filepath.Join(rootDir, "AGENTS.md"), "Root AGENTS")
	writeFile(t, filepath.Join(rootDir, "internal/repository/AGENTS.md"), "Repo AGENTS")

	// Create ADR
	writeFile(t, filepath.Join(rootDir, "docs/adr/2026-04-06-test.md"), "# Test ADR\nThis mentions internal/repository")

	// Create coding-standards.md
	codingContent := `## Strict Lessons
### CI & Infrastructure
* Infrastructure Lesson
### UI & Frontend
* Frontend Lesson
`
	writeFile(t, filepath.Join(rootDir, ".agents/workflows/coding-standards.md"), codingContent)

	// Create invariants.json
	writeFile(t, filepath.Join(rootDir, ".agents/invariants.json"), `{"db_engine": "sqlite"}`)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	if errClose := w.Close(); errClose != nil {
		t.Fatal(errClose)
	}
	os.Stdout = old

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Errorf("Capture failed: %v", err)
	}
	return output
}
