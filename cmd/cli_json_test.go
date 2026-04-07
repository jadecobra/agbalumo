package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestCLIJSONOutput(t *testing.T) {
	// Setup: Ensure we use a test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "cli_test.db")
	_ = os.Setenv("DATABASE_URL", dbPath)
	defer func() { _ = os.Unsetenv("DATABASE_URL") }()

	// 1. Test listing list --json (empty)
	t.Run("listing list --json empty", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"listing", "list"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := strings.TrimSpace(buf.String())
		// Cobra might print extra info, so we look for the JSON part
		if !strings.Contains(output, "[]") {
			t.Errorf("Expected output to contain '[]', got %q", output)
		}
	})

	// 2. Test category list --json
	t.Run("category list --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"category", "list"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		jsonPart := extractJSONFromOutput(t, buf.String())
		var categories []domain.CategoryData
		if err := json.Unmarshal([]byte(jsonPart), &categories); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
	})

	// 3. Test listing create --json
	t.Run("listing create --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"listing", "create", "--title", "JSON Test Listing"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		jsonPart := extractJSONFromOutput(t, buf.String())
		var listing domain.Listing
		if err := json.Unmarshal([]byte(jsonPart), &listing); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if listing.Title != "JSON Test Listing" {
			t.Errorf("Expected title 'JSON Test Listing', got %q", listing.Title)
		}
	})
}

func extractJSONFromOutput(t *testing.T, output string) string {
	start := strings.IndexAny(output, "[{")
	if start == -1 {
		t.Fatalf("No JSON found in output: %q", output)
	}
	return output[start:]
}
