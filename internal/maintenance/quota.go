package maintenance

import (
	"fmt"
	"os"
	"strings"
)

// QuotaConfig defines the paths and thresholds for quota enforcement.
type QuotaConfig struct {
	AgentsPath    string
	ResolverPath  string
	StandardsPath string
	ManifestPath  string
	MaxTaxBytes   int64
}

// CheckPreflightTax verifies the combined size of the preflight bundle.
func CheckPreflightTax(cfg QuotaConfig) error {
	paths := []string{
		cfg.AgentsPath,
		cfg.ResolverPath,
		cfg.StandardsPath,
		cfg.ManifestPath,
	}

	var totalSize int64
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		totalSize += info.Size()
	}

	if totalSize > cfg.MaxTaxBytes {
		return fmt.Errorf("preflight tax exceeded: combined size is %d bytes (limit: %d bytes). Prune AGENTS.md or standards to restore performance", totalSize, cfg.MaxTaxBytes)
	}

	return nil
}

// CheckQuotaViolation verifies if a commit message indicates an unauthorized high-tier model usage.
func CheckQuotaViolation(message string) error {
	expensiveMarkers := []string{"[Opus]", "[Pro]", "[Gemini 3.1 Pro]", "[3.1 Pro]"}

	isExpensive := false
	for _, marker := range expensiveMarkers {
		if containsIgnoreCase(message, marker) {
			isExpensive = true
			break
		}
	}

	if isExpensive && !containsIgnoreCase(message, "OVERRIDE") {
		return fmt.Errorf("QUOTA VIOLATION: High-tier model usage detected in commit message without OVERRIDE flag. Use Gemini 3 Flash for product execution or add OVERRIDE to the commit message")
	}

	return nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
