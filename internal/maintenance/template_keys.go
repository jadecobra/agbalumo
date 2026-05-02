package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckTemplateKeyGaps scans all templates to ensure that when a template is called via 'dict',
// all keys referenced within that template (e.g. $.Foo or .Foo) are explicitly passed.
func CheckTemplateKeyGaps(dir string) ([]DesignViolation, error) {
	registry, err := buildTemplateRegistry(dir)
	if err != nil {
		return nil, err
	}
	return scanTemplateInvocations(dir, registry)
}

func buildTemplateRegistry(dir string) (map[string]map[string]bool, error) {
	registry := make(map[string]map[string]bool)
	defineRegex := regexp.MustCompile(`{{\s*define\s*"([^"]+)"\s*}}`)

	// #nosec G122 -- local maintenance utility
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		// #nosec G304,G122 -- local templates only
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sContent := string(content)
		matches := defineRegex.FindAllStringSubmatchIndex(sContent, -1)

		for i, match := range matches {
			name := sContent[match[2]:match[3]]
			end := len(sContent)
			if i+1 < len(matches) {
				end = matches[i+1][0]
			}
			registry[name] = extractReferences(sContent[match[1]:end])
		}
		return nil
	})
	return registry, err
}

func extractReferences(block string) map[string]bool {
	refRegex := regexp.MustCompile(`(?:\$|\.)\.(\w+)`)
	refs := make(map[string]bool)
	refMatches := refRegex.FindAllStringSubmatch(block, -1)
	for _, rm := range refMatches {
		refs[rm[1]] = true
	}
	return refs
}

func scanTemplateInvocations(dir string, registry map[string]map[string]bool) ([]DesignViolation, error) {
	var violations []DesignViolation
	templateRegex := regexp.MustCompile(`{{\s*template\s*"([^"]+)"\s+dict\s+(.+?)\s*}}`)

	// #nosec G122 -- local maintenance utility
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		// #nosec G304,G122 -- local templates only
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sContent := string(content)
		lines := strings.Split(sContent, "\n")

		for lineIdx, line := range lines {
			matches := templateRegex.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				v := checkInvocation(path, lineIdx+1, line, m[1], m[2], registry)
				violations = append(violations, v...)
			}
		}
		return nil
	})
	return violations, err
}

func checkInvocation(path string, lineNum int, line, tmplName, dictContent string, registry map[string]map[string]bool) []DesignViolation {
	var violations []DesignViolation
	requiredKeys, exists := registry[tmplName]
	if !exists {
		return nil
	}

	keyRegex := regexp.MustCompile(`"(\w+)"`)
	keyMatches := keyRegex.FindAllStringSubmatch(dictContent, -1)
	passedKeys := make(map[string]bool)
	for _, km := range keyMatches {
		passedKeys[km[1]] = true
	}

	for req := range requiredKeys {
		if !passedKeys[req] {
			violations = append(violations, DesignViolation{
				File:    path,
				Line:    lineNum,
				Content: strings.TrimSpace(line),
				Reason:  fmt.Sprintf("Template '%s' references variable '%s' but it is not passed in dict", tmplName, req),
			})
		}
	}
	return violations
}
