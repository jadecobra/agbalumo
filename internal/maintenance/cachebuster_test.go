package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyCacheBuster(t *testing.T) {
	tmpDir := t.TempDir()

	cssDir := filepath.Join(tmpDir, "ui", "static", "css")
	htmlDir := filepath.Join(tmpDir, "ui", "templates", "components")

	if err := os.MkdirAll(cssDir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(htmlDir, 0750); err != nil {
		t.Fatal(err)
	}

	cssPath := filepath.Join(cssDir, "output.css")
	htmlPath := filepath.Join(htmlDir, "head_meta.html")

	// Write mock CSS
	if err := os.WriteFile(cssPath, []byte("body { color: red; }"), 0600); err != nil {
		t.Fatal(err)
	}
	// The SHA256 of "body { color: red; }" first 8 chars: 5de625c3

	// Test 1: Valid Hash
	htmlContent := `<link href="/static/css/output.css?v=5de625c3" rel="stylesheet" />`
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err != nil {
		t.Fatalf("expected nil, got err: %v", err)
	}

	// Test 2: Invalid Hash
	htmlContentInvalid := `<link href="/static/css/output.css?v=5" rel="stylesheet" />`
	if err := os.WriteFile(htmlPath, []byte(htmlContentInvalid), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err == nil {
		t.Fatalf("expected error for mismatched hash, got nil")
	}

	// Test 3: Missing pattern
	htmlContentMissing := `<link href="/static/css/output.css" rel="stylesheet" />`
	if err := os.WriteFile(htmlPath, []byte(htmlContentMissing), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err == nil {
		t.Fatalf("expected error for missing pattern, got nil")
	}
}
