package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyCacheBuster(t *testing.T) {
	tmpDir := t.TempDir()
	cssPath, htmlPath := setupTestCacheBusterPaths(t, tmpDir)

	// Write mock CSS
	if err := os.WriteFile(cssPath, []byte("body { color: red; }"), 0600); err != nil {
		t.Fatal(err)
	}

	// Test 1: Valid Hash & Versioned scripts
	t.Run("Valid Hash and Scripts", func(t *testing.T) {
		htmlContent := `<link href="/static/css/output.css?v=5de625c3" rel="stylesheet" />
<script src="/static/js/near-me.js?v=2" defer></script>
<script src="/static/js/htmx.min.js"></script>`
		writeTestHTML(t, htmlPath, htmlContent)
		if err := VerifyCacheBuster(tmpDir); err != nil {
			t.Fatalf("expected nil, got err: %v", err)
		}
	})

	// Test 2: Invalid CSS Hash
	t.Run("Invalid CSS Hash", func(t *testing.T) {
		htmlContent := `<link href="/static/css/output.css?v=5" rel="stylesheet" />
<script src="/static/js/near-me.js?v=2" defer></script>`
		writeTestHTML(t, htmlPath, htmlContent)
		if err := VerifyCacheBuster(tmpDir); err == nil {
			t.Fatal("expected error for mismatched hash")
		}
	})

	// Test 3: Unversioned JS script
	t.Run("Unversioned JS Script", func(t *testing.T) {
		htmlContent := `<link href="/static/css/output.css?v=5de625c3" rel="stylesheet" />
<script src="/static/js/near-me.js" defer></script>`
		writeTestHTML(t, htmlPath, htmlContent)
		if err := VerifyCacheBuster(tmpDir); err == nil {
			t.Fatal("expected error for unversioned custom JS script")
		}
	})
}

func setupTestCacheBusterPaths(t *testing.T, tmpDir string) (string, string) {
	t.Helper()
	cssDir := filepath.Join(tmpDir, "ui", "static", "css")
	htmlDir := filepath.Join(tmpDir, "ui", "templates", "components")

	if err := os.MkdirAll(cssDir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(htmlDir, 0750); err != nil {
		t.Fatal(err)
	}

	return filepath.Join(cssDir, "output.css"), filepath.Join(htmlDir, "head_meta.html")
}

func writeTestHTML(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}
