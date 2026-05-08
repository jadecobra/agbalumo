package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveIntent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal manifest fixture
	manifestContent := `
commands:
  - name: test
    trigger: test_authoring, test
    description: Run tests
  - name: browser
    trigger: ui_change
    description: UI tests
  - name: ci-tools
    trigger: ci_config_change
    description: CI Tools

skills:
  - name: go-tdd
    trigger: test_authoring, writing tests, feature_implementation
    path: .agents/skills/go-tdd/SKILL.md
  - name: browser-verify
    trigger: ui_change
    path: .agents/skills/browser-verify/SKILL.md
`
	agentsDir := filepath.Join(tmpDir, ".agents")
	err := os.Mkdir(agentsDir, 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(agentsDir, "verify-manifest.yaml"), []byte(manifestContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		intent   string
		expected []string // Name of matches
	}{
		{
			intent:   "writing tests",
			expected: []string{"go-tdd", "test"},
		},
		{
			intent:   "ui change",
			expected: []string{"browser", "browser-verify"}, // Should NOT match ci-tools
		},
		{
			intent:   "config change",
			expected: []string{"ci-tools"}, // Matches 2/3 words of ci_config_change
		},
		{
			intent:   "nonexistent intent",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.intent, func(t *testing.T) {
			matches, err := ResolveIntent(tmpDir, tt.intent)
			if err != nil {
				t.Fatalf("ResolveIntent failed: %v", err)
			}
			assertMatches(t, matches, tt.expected)
		})
	}
}

func assertMatches(t *testing.T, matches []ResolvedMatch, expected []string) {
	t.Helper()
	if len(matches) != len(expected) {
		t.Errorf("expected %d matches, got %d", len(expected), len(matches))
	}

	for _, exp := range expected {
		found := false
		for _, m := range matches {
			if m.Name == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected match %q not found", exp)
		}
	}
}
