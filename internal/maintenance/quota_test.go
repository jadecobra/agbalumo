package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPreflightTax(t *testing.T) {
	tmpDir := t.TempDir()

	// Create dummy preflight files
	agentsPath := filepath.Join(tmpDir, "AGENTS.md")
	resolverPath := filepath.Join(tmpDir, "RESOLVER.md")
	standardsPath := filepath.Join(tmpDir, "coding-standards.md")
	manifestPath := filepath.Join(tmpDir, "verify-manifest.yaml")

	files := []string{agentsPath, resolverPath, standardsPath, manifestPath}
	for _, f := range files {
		if err := os.WriteFile(f, make([]byte, 5000), 0600); err != nil {
			t.Fatal(err)
		}
	}

	// Total size = 20,000 bytes

	t.Run("Fails if above threshold", func(t *testing.T) {
		config := QuotaConfig{
			AgentsPath:    agentsPath,
			ResolverPath:  resolverPath,
			StandardsPath: standardsPath,
			ManifestPath:  manifestPath,
			MaxTaxBytes:   15000,
		}

		err := CheckPreflightTax(config)
		if err == nil {
			t.Error("expected error for exceeding quota, got nil")
		}
	})

	t.Run("Passes if below threshold", func(t *testing.T) {
		config := QuotaConfig{
			AgentsPath:    agentsPath,
			ResolverPath:  resolverPath,
			StandardsPath: standardsPath,
			ManifestPath:  manifestPath,
			MaxTaxBytes:   25000,
		}

		err := CheckPreflightTax(config)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestQuotaGate(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		expectErr bool
	}{
		{"Normal flash commit", "feat(auth): add login", false},
		{"Opus marker without override", "feat(auth): [Opus] refactor", true},
		{"Pro marker without override", "feat(auth): [Gemini 3.1 Pro] fix bug", true},
		{"Expensive with override", "feat(auth): [Opus] refactor OVERRIDE", false},
		{"Case insensitive override", "feat(auth): [Opus] refactor override", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckQuotaViolation(tt.message)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
