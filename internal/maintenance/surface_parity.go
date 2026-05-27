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

var forbiddenStyles = []struct {
	pattern     string
	description string
}{
	{"border-2 border-stone-", "neobrutalist thick outline (border-2 border-stone-*)"},
	{"shadow-[2px_2px_", "neobrutalist hard offset shadow"},
	{"shadow-[3px_3px_", "neobrutalist hard offset shadow"},
	{"shadow-[4px_4px_", "neobrutalist hard offset shadow"},
	{"shadow-[5px_5px_", "neobrutalist hard offset shadow"},
}

// CheckSurfaceParity ensures that key UI tokens are consistent between listing cards and modal details.
func CheckSurfaceParity(rootDir string) ([]string, error) {
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

	var violations []string
	violations = append(violations, checkTokenParity(sCard, sModal)...)
	violations = append(violations, checkForbiddenStyles(sCard, sModal)...)

	return violations, nil
}

// checkTokenParity checks required visual token consistency.
func checkTokenParity(card, modal string) []string {
	var violations []string
	for _, token := range requiredTokens {
		hasCard := strings.Contains(card, token)
		hasModal := strings.Contains(modal, token)

		if hasCard && !hasModal {
			violations = append(violations, fmt.Sprintf("Token '%s' present in listing_card.html but missing in modal_detail.html", token))
		}
		if !hasCard && hasModal {
			violations = append(violations, fmt.Sprintf("Token '%s' present in modal_detail.html but missing in listing_card.html", token))
		}
	}
	return violations
}

// checkForbiddenStyles checks for disallowed neobrutalist styling.
func checkForbiddenStyles(card, modal string) []string {
	var violations []string
	for _, fs := range forbiddenStyles {
		if strings.Contains(card, fs.pattern) {
			violations = append(violations, fmt.Sprintf("Forbidden style '%s' detected in listing_card.html: %s", fs.pattern, fs.description))
		}
		if strings.Contains(modal, fs.pattern) {
			violations = append(violations, fmt.Sprintf("Forbidden style '%s' detected in modal_detail.html: %s", fs.pattern, fs.description))
		}
	}
	return violations
}
