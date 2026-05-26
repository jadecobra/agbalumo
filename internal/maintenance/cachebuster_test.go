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

	// Test 1: Valid Hash & Versioned scripts
	htmlContent := `<link href="/static/css/output.css?v=5de625c3" rel="stylesheet" />
<script src="/static/js/near-me.js?v=2" defer></script>
<script src="/static/js/htmx.min.js"></script>`
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err != nil {
		t.Fatalf("expected nil, got err: %v", err)
	}

	// Test 2: Invalid CSS Hash but valid JS scripts
	htmlContentInvalid := `<link href="/static/css/output.css?v=5" rel="stylesheet" />
<script src="/static/js/near-me.js?v=2" defer></script>`
	if err := os.WriteFile(htmlPath, []byte(htmlContentInvalid), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err == nil {
		t.Fatalf("expected error for mismatched hash, got nil")
	}

	// Test 3: Missing CSS pattern but valid JS scripts
	htmlContentMissing := `<link href="/static/css/output.css" rel="stylesheet" />
<script src="/static/js/near-me.js?v=2" defer></script>`
	if err := os.WriteFile(htmlPath, []byte(htmlContentMissing), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err == nil {
		t.Fatalf("expected error for missing pattern, got nil")
	}

	// Test 4: Valid CSS Hash but unversioned custom JS script
	htmlContentUnversionedJS := `<link href="/static/css/output.css?v=5de625c3" rel="stylesheet" />
<script src="/static/js/near-me.js" defer></script>`
	if err := os.WriteFile(htmlPath, []byte(htmlContentUnversionedJS), 0600); err != nil {
		t.Fatal(err)
	}

	if err := VerifyCacheBuster(tmpDir); err == nil {
		t.Fatalf("expected error for unversioned custom JS script, got nil")
	}
}
