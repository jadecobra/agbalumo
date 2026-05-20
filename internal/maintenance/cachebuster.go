package maintenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// VerifyCacheBuster ensures that the CSS cache buster in head_meta.html
// matches the actual SHA256 hash of the output.css file.
func VerifyCacheBuster(rootDir string) error {
	fmt.Println("🔍 Verifying CSS cache buster hash...")

	cssPath := filepath.Join(rootDir, "ui", "static", "css", "output.css")
	htmlPath := filepath.Join(rootDir, "ui", "templates", "components", "head_meta.html")

	// 1. Compute hash of output.css
	// #nosec G304 // Reading static asset for hashing
	f, err := os.Open(cssPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️  output.css not found, skipping cache buster check.")
			return nil
		}
		return fmt.Errorf("failed to open output.css: %w", err)
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, copyErr := io.Copy(h, f); copyErr != nil {
		return fmt.Errorf("failed to hash output.css: %w", copyErr)
	}
	hashBytes := h.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)[:8] // Use first 8 chars for brevity

	// 2. Read head_meta.html
	// #nosec G304 // Reading template file for cache buster verification
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to read head_meta.html: %w", err)
	}

	// 3. Extract the current cache buster
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
	return nil
}
