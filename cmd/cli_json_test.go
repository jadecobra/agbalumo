package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestCLIJSONOutput(t *testing.T) {
	// Setup: Ensure we use a test database
	os.Setenv("DATABASE_URL", "@tester/cli_test.db")
	defer os.Unsetenv("DATABASE_URL")
	defer os.Remove("@tester/cli_test.db")

	// 1. Test listing list --json (empty)
	t.Run("listing list --json empty", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"listing", "list", "--json"})

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
		rootCmd.SetArgs([]string{"category", "list", "--json"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()
		// Find the JSON part (starts with [ or {)
		start := strings.IndexAny(output, "[{")
		if start == -1 {
			t.Fatalf("No JSON found in output: %q", output)
		}
		jsonPart := output[start:]

		var categories []domain.CategoryData
		if err := json.Unmarshal([]byte(jsonPart), &categories); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, jsonPart)
		}
	})

	// 3. Test listing create --json
	t.Run("listing create --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"listing", "create", "--title", "JSON Test Listing", "--json"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()
		start := strings.IndexAny(output, "[{")
		if start == -1 {
			t.Fatalf("No JSON found in output: %q", output)
		}
		jsonPart := output[start:]

		var listing domain.Listing
		if err := json.Unmarshal([]byte(jsonPart), &listing); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, jsonPart)
		}

		if listing.Title != "JSON Test Listing" {
			t.Errorf("Expected title 'JSON Test Listing', got %q", listing.Title)
		}
	})
}
