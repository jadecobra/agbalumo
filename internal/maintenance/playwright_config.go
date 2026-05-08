package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// VerifyPlaywrightConfig ensures that Playwright is configured to never open the HTML reporter automatically,
// which prevents port exhaustion from orphaned processes.
func VerifyPlaywrightConfig(dir string) error {
	path := filepath.Join(dir, "playwright.config.ts")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil // Not a playwright project or config missing, skip
	}

	// #nosec G304 // maintenance utility runs on trusted local directories
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read playwright config: %w", err)
	}

	// We look for the specific configuration that prevents automatic opening
	if !strings.Contains(string(content), "open: 'never'") {
		return fmt.Errorf("playwright config must contain `open: 'never'` in the html reporter options to prevent port exhaustion")
	}

	return nil
}
