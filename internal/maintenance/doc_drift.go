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

// CheckCommandDrift scans all documentation for stale command references.
func CheckCommandDrift(rootDir string) ([]DriftViolation, error) {
	docs, err := findMarkdownDocs(rootDir)
	if err != nil {
		return nil, err
	}

	subcommands := []string{
		"browser", "ci", "design", "doc-drift", "api-spec", "template-drift",
		"location-backfill", "enrich", "context-cost", "coverage", "audit", "verify-shas",
		"ci-tools", "js-syntax", "gitleaks", "ignored-files", "critique", "heal",
		"perf", "check-gates", "watch", "gosec-rationale", "preflight",
		"session-context", "janitor", "dump-invariants", "visual-audit",
		"skill-conformance", "check-resolvable", "map", "schema", "trace",
		"root-hygiene", "precommit", "test",
	}
	cmdMap := make(map[string]bool)
	for _, c := range subcommands {
		cmdMap[c] = true
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

	docs = append(docs, findExtraDocs(rootDir)...)
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

func findExtraDocs(rootDir string) []string {
	var extra []string
	stdPath := filepath.Join(rootDir, ".agents/workflows/coding-standards.md")
	if _, err := os.Stat(stdPath); err == nil {
		extra = append(extra, ".agents/workflows/coding-standards.md")
	}
	return extra
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
			if strings.HasPrefix(ref, ".agents/") {
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
