package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var requiredTokens = []string{
	"bg-white dark:bg-surface-dark",
	"border-stone-200 dark:border-stone-800",
	"text-text-main dark:text-earth-cream",
}

// CheckSurfaceParity ensures that key UI tokens are consistent between listing cards and modal details.
func CheckSurfaceParity(rootDir string) ([]string, error) {
	var violations []string

	cardPath := filepath.Join(rootDir, "ui/templates/partials/listing_card.html")
	modalPath := filepath.Join(rootDir, "ui/templates/partials/modal_detail.html")

	// #nosec G304 -- maintenance utility reads local templates
	cardContent, err := os.ReadFile(cardPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read listing_card.html: %w", err)
	}
	// #nosec G304 -- maintenance utility reads local templates
	modalContent, err := os.ReadFile(modalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read modal_detail.html: %w", err)
	}

	sCard := string(cardContent)
	sModal := string(modalContent)

	for _, token := range requiredTokens {
		hasCard := strings.Contains(sCard, token)
		hasModal := strings.Contains(sModal, token)

		if hasCard && !hasModal {
			violations = append(violations, fmt.Sprintf("Token '%s' present in listing_card.html but missing in modal_detail.html", token))
		}
		if !hasCard && hasModal {
			violations = append(violations, fmt.Sprintf("Token '%s' present in modal_detail.html but missing in listing_card.html", token))
		}
	}

	return violations, nil
}
