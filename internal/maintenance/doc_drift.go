package maintenance

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DriftViolation represents a stale file reference in documentation.
type DriftViolation struct {
	DocFile        string
	ReferencedPath string
	Line           int
	Exists         bool
}

var (
	backtickRegex = regexp.MustCompile("`([^`]+)`")
	boldRegex     = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// CheckDocDrift scans core documentation for stale file path references.
func CheckDocDrift(rootDir string) ([]DriftViolation, error) {
	docs := []string{"docs/architecture-overview.md", "docs/ATLAS.md"}
	var violations []DriftViolation
	for _, docRelPath := range docs {
		v, err := checkDocFile(rootDir, docRelPath)
		if err != nil {
			return nil, err
		}
		violations = append(violations, v...)
	}
	return violations, nil
}

func checkDocFile(rootDir, docRelPath string) ([]DriftViolation, error) {
	docPath := filepath.Join(rootDir, docRelPath)
	file, err := os.Open(filepath.Clean(docPath)) //nolint:gosec // G304: maintenance utility
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close() //nolint:errcheck // read-only check

	var violations []DriftViolation
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		v := checkLine(rootDir, docRelPath, lineNum, scanner.Text())
		violations = append(violations, v...)
	}
	return violations, nil
}

func checkLine(rootDir, docRelPath string, lineNum int, line string) []DriftViolation {
	var violations []DriftViolation
	matches := extractPotentialPaths(line)
	for _, ref := range matches {
		cleanRef := strings.TrimSpace(ref)
		if cleanRef == "" || !isLikelyPath(cleanRef) {
			continue
		}
		fullPath := filepath.Join(rootDir, cleanRef)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			violations = append(violations, DriftViolation{
				DocFile:        docRelPath,
				ReferencedPath: cleanRef,
				Line:           lineNum,
				Exists:         false,
			})
		}
	}
	return violations
}

func extractPotentialPaths(line string) []string {
	var matches []string
	matches = append(matches, extractByRegex(line, backtickRegex)...)
	matches = append(matches, extractByRegex(line, boldRegex)...)
	return matches
}

func extractByRegex(line string, re *regexp.Regexp) []string {
	var results []string
	matches := re.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) > 1 {
			p := cleanRefPath(m[1])
			if p != "" {
				results = append(results, p)
			}
		}
	}
	return results
}

func cleanRefPath(p string) string {
	// Handle func references like path/to/file.go::FuncName
	if strings.Contains(p, "::") {
		p = strings.Split(p, "::")[0]
	}
	// Skip commands (contain spaces) or URLs
	if strings.Contains(p, " ") || strings.HasPrefix(p, "http") {
		return ""
	}
	return p
}

func isLikelyPath(s string) bool {
	// Skip if contains function call syntax
	if strings.Contains(s, "(") || strings.Contains(s, ")") {
		return false
	}
	// Common top-level dirs/files in this project
	roots := []string{"internal", "cmd", "docs", "pkg", "api", "web", "assets", "tools", "scripts", "Dockerfile", "Makefile", "go.mod", "ui"}
	for _, r := range roots {
		if s == r || strings.HasPrefix(s, r+"/") {
			return true
		}
	}
	return false
}
