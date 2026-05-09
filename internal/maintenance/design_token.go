package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckTokenStrictness scans all .html templates for arbitrary Tailwind value brackets.
// It enforces the use of predefined design tokens for core layout and styling categories.
func CheckTokenStrictness(dir string) ([]DesignViolation, error) {
	var violations []DesignViolation

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		v, err := checkFileTokenStrictness(path)
		if err != nil {
			return err
		}
		violations = append(violations, v...)
		return nil
	})

	return violations, err
}

func checkFileTokenStrictness(path string) ([]DesignViolation, error) {
	var violations []DesignViolation

	// #nosec G304 -- utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	// Enforce predefined tokens for text, background, border, padding, margin, width, and height.
	// Rejects arbitrary values like text-[#123], p-[17px], w-[32.5%].
	tokenRegex := regexp.MustCompile(`\b(text|bg|border|p|m|w|h)-\[[^\]]+\]`)

	// Explicit whitelist for pervasive values not yet tokenized in tailwind.config.js
	// or values used for Material Symbols and architectural layout patterns.
	whitelist := map[string]bool{
		"text-[10px]":         true, // Caption size
		"text-[11px]":         true, // Nav item size
		"text-[12px]":         true, // Small icon size
		"text-[14px]":         true, // Standard icon size / small text
		"text-[16px]":         true, // Icon size
		"text-[18px]":         true, // Icon size
		"text-[20px]":         true, // Icon size
		"text-[24px]":         true, // Icon size
		"text-[28px]":         true, // Large icon size
		"text-[32px]":         true, // Hero icon size
		"text-[48px]":         true, // Modal icon size
		"h-[35vh]":            true, // Modal header aspect
		"h-[80px]":            true, // Header height
		"h-[85vh]":            true, // Modal max height
		"h-[90vh]":            true, // Modal max height
		"h-[90dvh]":           true, // Dynamic viewport height
		"h-[120px]":           true, // Header height
		"h-[280px]":           true, // Card height
		"w-[150px]":           true, // Table cell width
		"w-[200px]":           true, // Table cell width
		"w-[240px]":           true, // Featured card width
		"w-[400px]":           true, // Modal side drawer width
		"w-[calc(100%-2rem)]": true, // Floating mobile nav
		"border-[3px]":        true, // Brutalist border weight
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip lines with explicit override comment
		if strings.Contains(line, "verify:allow-arbitrary") {
			continue
		}

		matches := tokenRegex.FindAllString(line, -1)
		for _, match := range matches {
			if whitelist[match] {
				continue
			}
			violations = append(violations, DesignViolation{
				File:    path,
				Line:    lineNumber,
				Content: strings.TrimSpace(line),
				Reason:  fmt.Sprintf("Forbidden arbitrary Tailwind value '%s' (use predefined design tokens)", match),
			})
		}
	}

	return violations, scanner.Err()
}
