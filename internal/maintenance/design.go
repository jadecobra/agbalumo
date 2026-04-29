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

