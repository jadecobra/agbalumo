package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DesignViolation represents a deviation from the design system.
type DesignViolation struct {
	Content string
	Reason  string
	File    string
	Line    int
}

// CheckDesignStandards scans templates for violations of the UI Dialect protocol.
func CheckDesignStandards(dir string) ([]DesignViolation, error) {
	var violations []DesignViolation

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		v, err := checkFileStandards(path)
		if err != nil {
			return err
		}
		violations = append(violations, v...)

		uv, err := checkUppercaseDensity(path)
		if err != nil {
			return err
		}
		violations = append(violations, uv...)

		return nil
	})

	return violations, err
}

func checkFileStandards(path string) ([]DesignViolation, error) {
	var violations []DesignViolation

	// #nosec G304 -- verification utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	// All rounding (except -full for pills) is now forbidden project-wide.
	roundedRegex := regexp.MustCompile(`\brounded-(sm|md|lg|xl|2xl|3xl)\b`)
	hexRegex := regexp.MustCompile(`(?i)#([0-9a-f]{3}|[0-9a-f]{6})\b`)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		violations = append(violations, checkRounding(path, lineNumber, line, roundedRegex)...)
		violations = append(violations, checkHexCodes(path, lineNumber, line, hexRegex)...)
		violations = append(violations, checkMinFontSize(path, lineNumber, line)...)
		violations = append(violations, checkLowContrastOpacity(path, lineNumber, line)...)
		violations = append(violations, checkHardcodedModalBg(path, lineNumber, line)...)
		violations = append(violations, checkInlineHandlers(path, lineNumber, line)...)
		violations = append(violations, checkInlineStyles(path, lineNumber, line)...)
		violations = append(violations, checkMissingTestIDs(path, lineNumber, line)...)
		violations = append(violations, checkDeadSpacers(path, lineNumber, line)...)
		violations = append(violations, checkAdHocColors(path, lineNumber, line)...)
	}

	// File-level checks
	fv, err := checkRedundantCloseButtons(path)
	if err == nil {
		violations = append(violations, fv...)
	}

	return violations, scanner.Err()
}

func checkRounding(path string, lineNum int, line string, re *regexp.Regexp) []DesignViolation {
	var v []DesignViolation
	if match := re.FindString(line); match != "" {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  fmt.Sprintf("Forbidden rounding class '%s' (Brutalist standard requires sharp edges everywhere)", match),
		})
	}
	return v
}

func checkHexCodes(path string, lineNum int, line string, re *regexp.Regexp) []DesignViolation {
	var v []DesignViolation
	matches := re.FindAllStringIndex(line, -1)
	for _, matchIdx := range matches {
		start := matchIdx[0]
		// Skip HTML entities (&#123;).
		if start > 0 && line[start-1] == '&' {
			continue
		}
		match := line[matchIdx[0]:matchIdx[1]]
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  fmt.Sprintf("Hardcoded hex value '%s' (use Tailwind tokens instead)", match),
		})
	}
	return v
}

func checkMinFontSize(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	re := regexp.MustCompile(`text-\[(\d+)`)
	matches := re.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) > 1 {
			var size int
			_, err := fmt.Sscanf(m[1], "%d", &size)
			if err == nil && size < 10 {
				v = append(v, DesignViolation{
					File:    path,
					Line:    lineNum,
					Content: line,
					Reason:  "Font size below 10px minimum (ADR: surface-theme-unification)",
				})
				break
			}
		}
	}
	return v
}

func checkLowContrastOpacity(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	re := regexp.MustCompile(`text-text-sub/(\d+)`)
	matches := re.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) > 1 {
			var opacity int
			_, err := fmt.Sscanf(m[1], "%d", &opacity)
			if err == nil && opacity < 70 {
				v = append(v, DesignViolation{
					File:    path,
					Line:    lineNum,
					Content: line,
					Reason:  "Text opacity below 70% minimum for contrast",
				})
				break
			}
		}
	}
	return v
}

func checkHardcodedModalBg(path string, lineNum int, line string) []DesignViolation {
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "modal_") && base != "ui_components.html" {
		return nil
	}

	if strings.Contains(strings.ReplaceAll(line, "dark:bg-earth-dark", ""), "bg-earth-dark") {
		return []DesignViolation{{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Hardcoded dark background bypasses light/dark theme sync (ADR: surface-theme-unification)",
		}}
	}
	return nil
}

func checkInlineHandlers(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	re := regexp.MustCompile(`\bon(click|change|submit|mouseover)\s*=`)
	if re.MatchString(line) {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Inline event handler violates CSP policy (use data-* attributes + event delegation)",
		})
	}
	return v
}

func checkInlineStyles(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	if strings.Contains(line, "background-image:") {
		return nil
	}

	re := regexp.MustCompile(`\bstyle\s*=\s*"`)
	if re.MatchString(line) {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Inline style attribute (use Tailwind classes)",
		})
	}
	return v
}

func checkUppercaseDensity(path string) ([]DesignViolation, error) {
	// #nosec G304 -- verification utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	count := 0
	re := regexp.MustCompile(`\buppercase\b`)

	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			count++
		}
	}

	if count > 4 {
		return []DesignViolation{{
			File:   path,
			Line:   0, // File-level
			Reason: fmt.Sprintf("Uppercase density too high (%d occurrences, max 4 per template). Demote secondary elements to capitalize.", count),
		}}, scanner.Err()
	}

	return nil, scanner.Err()
}

func checkMissingTestIDs(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	// Target <button or <a that has hx-get|hx-post|hx-delete|hx-put
	if (strings.Contains(line, "<button") || strings.Contains(line, "<a ")) &&
		(strings.Contains(line, "hx-get") || strings.Contains(line, "hx-post") || strings.Contains(line, "hx-delete") || strings.Contains(line, "hx-put")) {
		if !strings.Contains(line, "data-testid=") {
			v = append(v, DesignViolation{
				File:    path,
				Line:    lineNum,
				Content: line,
				Reason:  "Interactive HTMX element missing data-testid (required for deterministic E2E testing)",
			})
		}
	}
	return v
}

func checkDeadSpacers(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	// Single-line check for empty section: <section...>\s*</section>
	re := regexp.MustCompile(`<section[^>]*>\s*</section>`)
	if re.MatchString(line) {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Empty section element creates dead vertical space",
		})
	}
	return v
}

func checkRedundantCloseButtons(path string) ([]DesignViolation, error) {
	// Only relevant for modal templates, excluding components definitions
	if !strings.Contains(path, "modal_") || strings.Contains(path, "ui_components.html") {
		return nil, nil
	}

	// #nosec G304 -- verification utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	closeActions := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "data-modal-action=\"close\"") || strings.Contains(line, "template \"btn_close\"") {
			closeActions++
		}
	}

	if closeActions > 1 {
		return []DesignViolation{{
			File:   path,
			Reason: fmt.Sprintf("Multiple close actions detected (%d). Modals MUST have exactly one primary close mechanism to prevent UI clutter.", closeActions),
		}}, nil
	}
	return nil, nil
}

func checkAdHocColors(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	// Catch things like bg-white/5 or hover:bg-black/5 which are often used for "fake" transparency that breaks themes
	re := regexp.MustCompile(`(bg|text)-(white|black|stone-\d+)/\d+`)
	matches := re.FindAllString(line, -1)
	for _, match := range matches {
		// Exceptions for backdrop or specific gradients if needed, but generally these should be themed tokens
		if strings.Contains(line, "backdrop:") || strings.Contains(line, "bg-gradient") {
			continue
		}
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  fmt.Sprintf("Ad-hoc color transparency '%s' detected. Use semantic tokens (e.g., bg-surface-dark/50) or theme-aware classes.", match),
		})
	}
	return v
}
