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

		v, err := checkFileStandards(path, info.Name())
		if err != nil {
			return err
		}
		violations = append(violations, v...)
		return nil
	})

	return violations, err
}

func checkFileStandards(path, filename string) ([]DesignViolation, error) {
	var violations []DesignViolation
	isSharpContext := strings.HasPrefix(filename, "admin_") || strings.Contains(path, "/admin/")

	// #nosec G304 -- verification utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	roundedRegex := regexp.MustCompile(`\brounded-(md|lg|xl|2xl|3xl|full)\b`)
	hexRegex := regexp.MustCompile(`(?i)#([0-9a-f]{3}|[0-9a-f]{6})\b`)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if isSharpContext {
			violations = append(violations, checkRounding(path, lineNumber, line, roundedRegex)...)
		}
		violations = append(violations, checkHexCodes(path, lineNumber, line, hexRegex)...)
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
			Reason:  fmt.Sprintf("Forbidden rounding class '%s' in Sharp (Admin) context", match),
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
