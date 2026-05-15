package maintenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// A11yViolationMap represents a mapped a11y violation.
type A11yViolationMap struct {
	ViolationID    string
	Impact         string
	TemplateFile   string
	HTMLSnippet    string
	FixSuggestion  string
	FailureSummary string
	Line           int
}

// MapA11yViolations parses the latest Playwright test results and maps violations to template files.
func MapA11yViolations(rootDir string) ([]A11yViolationMap, error) {
	latestFile, err := findLatestErrorContext(filepath.Join(rootDir, "test-results"))
	if err != nil {
		return nil, err
	}
	if latestFile == "" {
		return nil, nil
	}

	axeViolations, err := parseA11yViolations(latestFile)
	if err != nil {
		return nil, err
	}

	return mapViolationsToTemplates(rootDir, axeViolations), nil
}

func findLatestErrorContext(resultsDir string) (string, error) {
	var latestFile string
	var latestModTime int64

	err := filepath.Walk(resultsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Base(path) == "error-context.md" {
			if info.ModTime().Unix() > latestModTime {
				latestModTime = info.ModTime().Unix()
				latestFile = path
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to walk test-results: %w", err)
	}
	return latestFile, nil
}

type axeNode struct {
	HTML           string `json:"html"`
	FailureSummary string `json:"failureSummary"`
}

type axeViolation struct {
	ID     string    `json:"id"`
	Impact string    `json:"impact"`
	Nodes  []axeNode `json:"nodes"`
}

func parseA11yViolations(filePath string) ([]axeViolation, error) {
	content, err := os.ReadFile(filePath) //nolint:gosec // maintenance utility reads test results
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	marker := "A11y Violations: "
	idx := strings.Index(string(content), marker)
	if idx == -1 {
		return nil, nil
	}

	jsonStr := string(content)[idx+len(marker):]
	endIdx := strings.LastIndex(jsonStr, "]")
	if endIdx == -1 {
		return nil, fmt.Errorf("malformed a11y violations JSON in %s", filePath)
	}
	jsonStr = jsonStr[:endIdx+1]

	var violations []axeViolation
	if err := json.Unmarshal([]byte(jsonStr), &violations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal a11y violations: %w", err)
	}
	return violations, nil
}

func mapViolationsToTemplates(rootDir string, axeViolations []axeViolation) []A11yViolationMap {
	var results []A11yViolationMap
	templatesDir := filepath.Join(rootDir, "ui", "templates")

	for _, v := range axeViolations {
		for _, node := range v.Nodes {
			results = append(results, mapNodeToViolation(rootDir, templatesDir, v, node))
		}
	}
	return results
}

func mapNodeToViolation(rootDir, templatesDir string, v axeViolation, node axeNode) A11yViolationMap {
	mapping := A11yViolationMap{
		ViolationID:    v.ID,
		Impact:         v.Impact,
		HTMLSnippet:    node.HTML,
		FailureSummary: node.FailureSummary,
		FixSuggestion:  getFixSuggestion(v.ID),
		TemplateFile:   "unknown",
	}

	file, line := findInTemplates(templatesDir, node.HTML)
	if file != "" {
		if rel, err := filepath.Rel(rootDir, file); err == nil {
			mapping.TemplateFile = rel
		} else {
			mapping.TemplateFile = file
		}
		mapping.Line = line
	}
	return mapping
}

func getFixSuggestion(id string) string {
	switch id {
	case "label":
		return "Use <label for='ID'> or aria-label to associate text with the input."
	case "image-alt":
		return "Add a descriptive alt attribute: alt='description' or alt='' for decorative images."
	case "aria-required-attr":
		return "Add the missing required ARIA attributes for this role."
	default:
		return "Run `go run ./cmd/verify design` for static analysis."
	}
}

func findInTemplates(dir string, snippet string) (string, int) {
	var foundFile string
	var foundLine int
	target := strings.TrimSpace(snippet)

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		content, err := os.ReadFile(path) //nolint:gosec // maintenance utility reads template files
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, target) {
				foundFile = path
				foundLine = i + 1
				return filepath.SkipDir
			}
		}
		return nil
	})

	return foundFile, foundLine
}
