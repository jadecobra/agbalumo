package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// TemplateSchema represents the prop validation rules for a template component.
type TemplateSchema struct {
	Required  map[string]bool
	Optional  map[string]bool
	HasSchema bool
}

// CheckTemplateKeyGaps scans all templates to ensure that when a template is called via 'dict',
// all required keys are passed, and no extraneous keys are passed if the template has a schema.
func CheckTemplateKeyGaps(dir string) ([]DesignViolation, error) {
	registry, err := buildTemplateRegistry(dir)
	if err != nil {
		return nil, err
	}
	return scanTemplateInvocations(dir, registry)
}

func buildTemplateRegistry(dir string) (map[string]TemplateSchema, error) {
	registry := make(map[string]TemplateSchema)
	defineRegex := regexp.MustCompile(`{{\s*define\s*"([^"]+)"\s*}}`)
	propsRegex := regexp.MustCompile(`<!--\s*props:\s*\n([\s\S]*?)\n\s*-->`)

	// #nosec G122 -- local maintenance utility
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}
		return parseAndRegisterTemplates(path, registry, defineRegex, propsRegex)
	})
	return registry, err
}

func parseAndRegisterTemplates(path string, registry map[string]TemplateSchema, defineRegex, propsRegex *regexp.Regexp) error {
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
			if end < match[1] {
				end = len(sContent)
			}
		}
		block := sContent[match[1]:end]
		registry[name] = parseSingleTemplateSchema(block, propsRegex)
	}
	return nil
}

func parseSingleTemplateSchema(block string, propsRegex *regexp.Regexp) TemplateSchema {
	schema := TemplateSchema{
		Required:  make(map[string]bool),
		Optional:  make(map[string]bool),
		HasSchema: false,
	}

	propsMatch := propsRegex.FindStringSubmatch(block)
	if len(propsMatch) > 1 {
		schema.HasSchema = true
		parseYAMLProps(propsMatch[1], &schema)
	} else {
		// Fallback to legacy behavior: extract dynamic variables used in block
		refs := extractReferences(block)
		for req := range refs {
			schema.Required[req] = true
		}
	}
	return schema
}

func parseYAMLProps(commentBlock string, schema *TemplateSchema) {
	if strings.TrimSpace(commentBlock) == "" {
		return
	}
	var rawProps map[string]string
	err := yaml.Unmarshal([]byte(commentBlock), &rawProps)
	if err == nil {
		for k, v := range rawProps {
			cleanVal := strings.ToLower(strings.TrimSpace(v))
			if cleanVal == "required" {
				schema.Required[k] = true
			} else {
				schema.Optional[k] = true
			}
		}
	}
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

func scanTemplateInvocations(dir string, registry map[string]TemplateSchema) ([]DesignViolation, error) {
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

func parseDictKeys(dictContent string) []string {
	args := tokenizeDictArgs(dictContent)
	return extractKeysFromArgs(args)
}

func extractKeysFromArgs(args []string) []string {
	var keys []string
	for idx, arg := range args {
		if idx%2 == 0 {
			if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") && len(arg) >= 2 {
				keys = append(keys, arg[1:len(arg)-1])
			}
		}
	}
	return keys
}

func tokenizeDictArgs(dictContent string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	parenDepth := 0

	for i := 0; i < len(dictContent); i++ {
		ch := dictContent[i]
		inQuotes, parenDepth = updateState(ch, i, dictContent, inQuotes, parenDepth)
		args = appendToken(ch, inQuotes, parenDepth, &current, args)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func appendToken(ch byte, inQuotes bool, parenDepth int, current *strings.Builder, args []string) []string {
	if inQuotes || parenDepth > 0 || ch == ')' {
		current.WriteByte(ch)
		return args
	}

	if isWhitespace(ch) {
		if current.Len() > 0 {
			args = append(args, current.String())
			current.Reset()
		}
	} else {
		current.WriteByte(ch)
	}
	return args
}

func updateState(ch byte, i int, s string, inQuotes bool, parenDepth int) (bool, int) {
	if inQuotes {
		if ch == '"' && i > 0 && s[i-1] != '\\' {
			return false, parenDepth
		}
		return true, parenDepth
	}

	if ch == '"' {
		return true, parenDepth
	}

	if ch == '(' {
		return false, parenDepth + 1
	}

	if ch == ')' && parenDepth > 0 {
		return false, parenDepth - 1
	}

	return false, parenDepth
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func checkInvocation(path string, lineNum int, line, tmplName, dictContent string, registry map[string]TemplateSchema) []DesignViolation {
	schema, exists := registry[tmplName]
	if !exists {
		return nil
	}

	passedKeysList := parseDictKeys(dictContent)
	passedKeys := make(map[string]bool)
	for _, pk := range passedKeysList {
		passedKeys[pk] = true
	}

	var violations []DesignViolation
	violations = append(violations, checkRequiredKeys(path, lineNum, line, tmplName, schema, passedKeys)...)
	violations = append(violations, checkExtraneousKeys(path, lineNum, line, tmplName, schema, passedKeys)...)
	return violations
}

func checkRequiredKeys(path string, lineNum int, line, tmplName string, schema TemplateSchema, passedKeys map[string]bool) []DesignViolation {
	var violations []DesignViolation
	for req := range schema.Required {
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

func checkExtraneousKeys(path string, lineNum int, line, tmplName string, schema TemplateSchema, passedKeys map[string]bool) []DesignViolation {
	var violations []DesignViolation
	if !schema.HasSchema {
		return nil
	}
	for passed := range passedKeys {
		if !schema.Required[passed] && !schema.Optional[passed] {
			violations = append(violations, DesignViolation{
				File:    path,
				Line:    lineNum,
				Content: strings.TrimSpace(line),
				Reason:  fmt.Sprintf("Template '%s' accepts properties %v but got extraneous property '%s'", tmplName, getValidKeys(schema), passed),
			})
		}
	}
	return violations
}

func getValidKeys(schema TemplateSchema) []string {
	var keys []string
	for k := range schema.Required {
		keys = append(keys, k)
	}
	for k := range schema.Optional {
		keys = append(keys, k)
	}
	return keys
}
