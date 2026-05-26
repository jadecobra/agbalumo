package maintenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyCacheBuster ensures that:
// 1. The CSS cache buster in head_meta.html matches the actual SHA256 hash of output.css.
// 2. All custom local JavaScript assets (/static/js/*.js excluding .min.js) loaded in any template have ?v= query parameters to prevent stale browser caches.
//
// Replaces Strict Lesson: Client-Side JavaScript Cache Busting
func VerifyCacheBuster(rootDir string) error {
	fmt.Println("🔍 Verifying CSS cache buster hash...")

	cssPath := filepath.Join(rootDir, "ui", "static", "css", "output.css")
	htmlPath := filepath.Join(rootDir, "ui", "templates", "components", "head_meta.html")

	// 1. Compute hash of output.css
	// #nosec G304 - Reading static asset for hashing
	f, err := os.Open(cssPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️  output.css not found, skipping cache buster check.")
		} else {
			return fmt.Errorf("failed to open output.css: %w", err)
		}
	} else {
		defer func() { _ = f.Close() }()

		h := sha256.New()
		if _, copyErr := io.Copy(h, f); copyErr != nil {
			return fmt.Errorf("failed to hash output.css: %w", copyErr)
		}
		hashBytes := h.Sum(nil)
		hashStr := hex.EncodeToString(hashBytes)[:8] // Use first 8 chars for brevity

		// Read head_meta.html
		// #nosec G304 - Reading template file for cache buster verification
		htmlContent, readErr := os.ReadFile(htmlPath)
		if readErr != nil {
			return fmt.Errorf("failed to read head_meta.html: %w", readErr)
		}

		// Extract the current CSS cache buster
		re := regexp.MustCompile(`output\.css\?v=([a-fA-F0-9]+)`)
		matches := re.FindSubmatch(htmlContent)
		if len(matches) < 2 {
			return fmt.Errorf("failed to find 'output.css?v=...' pattern in head_meta.html")
		}

		currentHash := string(matches[1])

		if currentHash != hashStr {
			fmt.Printf("❌ Cache buster mismatch!\n   Expected hash: %s\n   Found hash:    %s\n", hashStr, currentHash)
			return fmt.Errorf("head_meta.html cache buster is out of date (must be %s)", hashStr)
		}

		fmt.Printf("✅ Cache buster is up to date (hash: %s).\n", hashStr)
	}

	// 2. Scan all HTML templates for unversioned local custom script tags
	fmt.Println("🔍 Verifying JavaScript cache busters in templates...")
	templatesDir := filepath.Join(rootDir, "ui", "templates")

	reScript := regexp.MustCompile(`(?i)<script\b[^>]*\bsrc\s*=\s*["']([^"']+)["']`)

	// #nosec G122 - Internal template walk is safe
	err = filepath.Walk(templatesDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// #nosec G304 - Internal template verification check
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("failed to read template %s: %w", path, readErr)
		}
		content := string(data)

		matches := reScript.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			src := match[1]

			// Only enforce on custom local scripts, ignore vendor minified files (.min.js)
			if strings.HasPrefix(src, "/static/js/") && !strings.Contains(src, ".min.js") {
				if !strings.Contains(src, "?v=") {
					relPath, _ := filepath.Rel(rootDir, path)
					return fmt.Errorf("unversioned custom script found in %s: src=%q (all custom JS scripts must have a '?v=' cache buster parameter)", relPath, src)
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("✅ All custom JavaScript script references are cache-busted.")
	return nil
}
