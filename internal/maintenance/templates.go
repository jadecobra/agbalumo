package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// ExtractRendererFunctions extracts defined functions from the renderer.
func ExtractRendererFunctions(path string) ([]string, error) {
	data, err := readFileOrErr(path, "renderer file")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var funcs []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Assuming format: "funcName": someFunc,
		if strings.HasPrefix(line, "\"") && strings.Contains(line, "\":") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				funcs = append(funcs, parts[1])
			}
		}
	}

	return uniqueStrings(funcs), nil
}

// ExtractTemplateFunctionCalls extracts used functions from template files.
func ExtractTemplateFunctionCalls(dir string) ([]string, error) {
	var used []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, domain.ExtHTML) {
			return err
		}

		data, err := readFileOrErr(path, "template file")
		if err != nil {
			return err
		}

		used = append(used, parseTemplateFunctions(string(data))...)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return uniqueStrings(used), nil
}

func parseTemplateFunctions(content string) []string {
	var used []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		used = append(used, extractFromLine(line)...)
	}
	return used
}

func extractFromLine(line string) []string {
	var used []string
	// Extract from {{ ... }}
	extract := func(sep string, stripRange bool) {
		if strings.Contains(line, sep) {
			parts := strings.Split(line, sep)
			for _, p := range parts[1:] {
				used = append(used, extractFirstWord(p, stripRange)...)
			}
		}
	}

	extract("{{", true)
	extract("|", false)
	return used
}

func extractFirstWord(input string, stripRange bool) []string {
	inner := strings.TrimSpace(input)
	if stripRange && strings.HasPrefix(inner, "range") {
		inner = strings.TrimSpace(strings.TrimPrefix(inner, "range"))
	}

	words := strings.FieldsFunc(inner, func(r rune) bool {
		return r == ' ' || r == '}' || r == '|' || r == '(' || r == ')'
	})

	if len(words) > 0 {
		name := words[0]
		if !isTemplateKeyword(name) && !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "$") {
			return []string{name}
		}
	}
	return nil
}

func isTemplateKeyword(s string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "end": true, "range": true, "with": true,
		"define": true, "block": true, "template": true, "nil": true, "len": true,
		"and": true, "or": true, "not": true, "index": true, "slice": true,
		"printf": true, "print": true, "println": true, "html": true,
		"urlquery": true, "js": true, "call": true,
	}
	return keywords[s]
}

// CheckTemplateDrift finds used functions that are not defined in the renderer.
func CheckTemplateDrift(defined, used []string) []string {
	var drifts []string
	defMap := make(map[string]bool)
	for _, d := range defined {
		defMap[d] = true
	}

	for _, u := range used {
		if !defMap[u] {
			// Skip capitalized words which are usually types or exports
			if len(u) > 0 && u[0] >= 'A' && u[0] <= 'Z' {
				continue
			}
			drifts = append(drifts, fmt.Sprintf("Undefined template function used: '%s'", u))
		}
	}
	return drifts
}
