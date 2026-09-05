package maintenance

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// DriftViolation represents a stale file reference in documentation.
type DriftViolation struct {
	DocFile        string
	ReferencedPath string
	Line           int
	Exists         bool
}

type manifest struct {
	Commands []struct {
		Name string `yaml:"name"`
	} `yaml:"commands"`
	Tools []struct {
		Name string `yaml:"name"`
	} `yaml:"tools"`
}

var (
	backtickRegex = regexp.MustCompile("`([^`]+)`")
	boldRegex     = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// CheckDocDrift scans all documentation in docs/ and specific .agents/ files for stale file path references.
func CheckDocDrift(rootDir string) ([]DriftViolation, error) {
	docs, err := findMarkdownDocs(rootDir)
	if err != nil {
		return nil, err
	}

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

func loadValidSubcommands(rootDir string) (map[string]bool, error) {
	manifestPath := filepath.Join(rootDir, ".agents/verify-manifest.yaml")
	data, err := os.ReadFile(manifestPath) //nolint:gosec // maintenance utility reads project-root manifest
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]bool), nil
		}
		return nil, err
	}

	var m manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	cmdMap := make(map[string]bool)
	for _, c := range m.Commands {
		if c.Name != "" {
			cmdMap[c.Name] = true
		}
	}
	for _, t := range m.Tools {
		if t.Name != "" {
			cmdMap[t.Name] = true
		}
	}
	return cmdMap, nil
}

// CheckCommandDrift scans all documentation for stale command references.
func CheckCommandDrift(rootDir string) ([]DriftViolation, error) {
	docs, err := findMarkdownDocs(rootDir)
	if err != nil {
		return nil, err
	}

	cmdMap, err := loadValidSubcommands(rootDir)
	if err != nil {
		return nil, err
	}

	var violations []DriftViolation
	for _, docRelPath := range docs {
		v, err := checkDocCommands(rootDir, docRelPath, cmdMap)
		if err != nil {
			return nil, err
		}
		violations = append(violations, v...)
	}
	return violations, nil
}

// CheckConfigPathDrift scans documentation for stale references to .agents/ config files.
func CheckConfigPathDrift(rootDir string) ([]DriftViolation, error) {
	docs, err := findMarkdownDocs(rootDir)
	if err != nil {
		return nil, err
	}

	var violations []DriftViolation
	for _, docRelPath := range docs {
		v, err := checkDocConfigPaths(rootDir, docRelPath)
		if err != nil {
			return nil, err
		}
		violations = append(violations, v...)
	}
	return violations, nil
}

func findMarkdownDocs(rootDir string) ([]string, error) {
	docsDir := filepath.Join(rootDir, "docs")
	docs, err := walkDocs(rootDir, docsDir)
	if err != nil {
		return nil, err
	}

	docs = append(docs, findAgentsAndDocFiles(rootDir)...)
	return docs, nil
}

func walkDocs(rootDir, docsDir string) ([]string, error) {
	var docs []string
	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		return handleWalk(rootDir, path, info, err, &docs)
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return docs, nil
}

func handleWalk(rootDir, path string, info os.FileInfo, err error, docs *[]string) error {
	if err != nil {
		return err
	}
	rel, skip := shouldSkipMarkdown(rootDir, path, info)
	if !skip {
		return nil
	}
	if info.IsDir() {
		return filepath.SkipDir
	}
	if rel != "" && strings.HasSuffix(rel, ".md") {
		*docs = append(*docs, rel)
	}
	return nil
}

func shouldSkipMarkdown(rootDir, path string, info os.FileInfo) (string, bool) {
	if info.IsDir() {
		return "", info.Name() == "adr" || info.Name() == "openapi"
	}
	rel, _ := filepath.Rel(rootDir, path)
	return rel, true
}

func findAgentsAndDocFiles(rootDir string) []string {
	var docs []string
	_ = filepath.Walk(rootDir, agentsWalkFn(rootDir, &docs))
	return docs
}

func agentsWalkFn(rootDir string, docs *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return handleAgentsDir(info.Name())
		}
		collectAgentDoc(rootDir, path, info, docs)
		return nil
	}
}

func handleAgentsDir(name string) error {
	if shouldSkipAgentsWalkDir(name) {
		return filepath.SkipDir
	}
	return nil
}

func collectAgentDoc(rootDir, path string, info os.FileInfo, docs *[]string) {
	if isTargetAgentDoc(rootDir, path, info) {
		if rel, err := filepath.Rel(rootDir, path); err == nil {
			*docs = append(*docs, rel)
		}
	}
}

func shouldSkipAgentsWalkDir(name string) bool {
	if name == "." || name == "" {
		return false
	}
	if strings.HasPrefix(name, ".") && name != ".agents" {
		return true
	}
	return name == vVendor || name == "node_modules" || name == "testdata" || name == ".tester"
}

func isTargetAgentDoc(rootDir, path string, info os.FileInfo) bool {
	if info.Name() == "AGENTS.md" {
		return true
	}
	rel, err := filepath.Rel(rootDir, path)
	if err != nil {
		return false
	}
	return strings.HasPrefix(rel, ".agents"+string(filepath.Separator)) && strings.HasSuffix(rel, ".md")
}

func checkDocCommands(rootDir, docRelPath string, validSubcmds map[string]bool) ([]DriftViolation, error) {
	return scanDoc(rootDir, docRelPath, func(line string) []string {
		return extractCommandsFromLine(line)
	}, func(ref string) bool {
		return validateCommand(ref, validSubcmds)
	})
}

func extractCommandsFromLine(line string) []string {
	var refs []string
	for _, m := range backtickRegex.FindAllStringSubmatch(line, -1) {
		if len(m) > 1 {
			cmd := strings.TrimSpace(m[1])
			if strings.HasPrefix(cmd, "task ") || strings.HasPrefix(cmd, "go run ./cmd/verify ") {
				refs = append(refs, cmd)
			}
		}
	}
	return refs
}

func validateCommand(ref string, validSubcmds map[string]bool) bool {
	if strings.HasPrefix(ref, "task ") {
		return false // task is stale
	}
	// Check if verify subcommand exists
	parts := strings.Fields(ref)
	if len(parts) >= 4 && parts[0] == "go" && parts[1] == "run" && parts[2] == "./cmd/verify" {
		sub := parts[3]
		return validSubcmds[sub]
	}
	return true
}

func checkDocConfigPaths(rootDir, docRelPath string) ([]DriftViolation, error) {
	return scanDoc(rootDir, docRelPath, func(line string) []string {
		return extractConfigPathsFromLine(line)
	}, func(ref string) bool {
		_, err := os.Stat(filepath.Join(rootDir, ref))
		return err == nil
	})
}

func extractConfigPathsFromLine(line string) []string {
	var refs []string
	for _, m := range backtickRegex.FindAllStringSubmatch(line, -1) {
		if len(m) > 1 {
			ref := strings.TrimSpace(m[1])
			if strings.HasPrefix(ref, ".agents/") && !strings.ContainsAny(ref, "*<>[]") {
				refs = append(refs, ref)
			}
		}
	}
	return refs
}

type lineExtractor func(line string) []string
type refValidator func(ref string) bool

func scanDoc(rootDir, docRelPath string, extract lineExtractor, validate refValidator) ([]DriftViolation, error) {
	docPath := filepath.Join(rootDir, docRelPath)
	file, err := os.Open(filepath.Clean(docPath))
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
		line := scanner.Text()
		refs := extract(line)
		for _, ref := range refs {
			if !validate(ref) {
				violations = append(violations, DriftViolation{
					DocFile:        docRelPath,
					ReferencedPath: ref,
					Line:           lineNum,
					Exists:         false,
				})
			}
		}
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
		fullPath := filepath.Join(rootDir, resolvePathToStat(cleanRef))
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

func resolvePathToStat(ref string) string {
	if ext := filepath.Ext(ref); len(ext) > 1 && ext[1] >= 'A' && ext[1] <= 'Z' {
		return strings.TrimSuffix(ref, ext)
	}
	return ref
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
	// Skip if contains function call syntax or wildcards/placeholders
	if strings.Contains(s, "(") || strings.Contains(s, ")") || strings.ContainsAny(s, "*<>[]") {
		return false
	}
	// Common top-level dirs/files in this project
	roots := []string{"internal", "cmd", "docs", "pkg", "api", "web", "assets", "tools", "scripts", "Dockerfile", "Makefile", "go.mod", "ui", ".agents"}
	for _, r := range roots {
		if s == r || strings.HasPrefix(s, r+"/") {
			return true
		}
	}
	return false
}
