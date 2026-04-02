package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

func VerifyTemplateDrift() bool {
	fmt.Println("Running Template Function Drift Check...")

	rendererPath := "internal/ui/renderer.go"
	templatesDir := "ui/templates"

	definedFuncs, err := ExtractRendererFunctions(rendererPath)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return false
	}

	usedFuncs, err := ExtractTemplateFunctionCalls(templatesDir)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return false
	}

	drifts := CheckTemplateDrift(definedFuncs, usedFuncs)
	if len(drifts) == 0 {
		fmt.Println("✅ All template functions are in sync.")
		return true
	}

	for _, d := range drifts {
		fmt.Printf("❌ %s\n", d)
	}
	fmt.Println("❌ Gate FAIL: Template Function Drift Detected!")
	return false
}

func ExtractRendererFunctions(path string) ([]string, error) {
	data, err := util.SafeReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read renderer file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var funcs []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"") && strings.Contains(line, "\":") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				funcs = append(funcs, parts[1])
			}
		}
	}

	return util.UniqueStrings(funcs), nil
}

func ExtractTemplateFunctionCalls(dir string) ([]string, error) {
	var used []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			data, err := util.SafeReadFile(path)
			if err != nil {
				return err
			}

			content := string(data)
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				if strings.Contains(line, "{{") {
					parts := strings.Split(line, "{{")
					for _, p := range parts[1:] {
						inner := strings.TrimSpace(p)
						if strings.HasPrefix(inner, "range") {
							inner = strings.TrimSpace(strings.TrimPrefix(inner, "range"))
						}

						words := strings.FieldsFunc(inner, func(r rune) bool {
							return r == ' ' || r == '}' || r == '|' || r == '(' || r == ')'
						})
						if len(words) > 0 {
							name := words[0]
							if !isTemplateKeyword(name) && !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "$") {
								used = append(used, name)
							}
						}
					}
				}

				if strings.Contains(line, "|") {
					parts := strings.Split(line, "|")
					for _, p := range parts[1:] {
						inner := strings.TrimSpace(p)
						words := strings.FieldsFunc(inner, func(r rune) bool {
							return r == ' ' || r == '}' || r == '|' || r == '(' || r == ')'
						})
						if len(words) > 0 {
							name := words[0]
							if !isTemplateKeyword(name) && !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "$") {
								used = append(used, name)
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return util.UniqueStrings(used), nil
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

func CheckTemplateDrift(defined, used []string) []string {
	var drifts []string
	defMap := make(map[string]bool)
	for _, d := range defined {
		defMap[d] = true
	}

	for _, u := range used {
		if !defMap[u] {
			if len(u) > 0 && u[0] >= 'A' && u[0] <= 'Z' {
				continue
			}
			drifts = append(drifts, fmt.Sprintf("Undefined template function used: '%s'", u))
		}
	}
	return drifts
}
