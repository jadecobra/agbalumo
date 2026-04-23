package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunSessionContext dumps all relevant rules, constraints, and ADRs for a specified path.
func RunSessionContext(rootDir, targetPath string) error {
	relPath, err := filepath.Rel(rootDir, targetPath)
	if err != nil {
		relPath = targetPath
	}
	relPath = filepath.Clean(relPath)

	fmt.Printf("📋 Session Context for: %s/\n", relPath)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	printSessionAGENTS(rootDir, relPath)
	printDomainLessons(rootDir, relPath)
	printRelatedADRs(rootDir, relPath)
	printInvariants(rootDir)

	return nil
}

func printSessionAGENTS(rootDir, relPath string) {
	current := relPath
	first := true
	for {
		agentsPath := filepath.Join(current, "AGENTS.md")
		fullPath := filepath.Join(rootDir, agentsPath)

		first = tryPrintAgentsFile(fullPath, agentsPath, first)

		if current == "." || current == "/" || current == "" {
			break
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
}

func tryPrintAgentsFile(fullPath, relPath string, isFirst bool) bool {
	// #nosec G304 -- path is derived from trusted repo structure
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return isFirst
	}

	label := "📁 Inherited AGENTS.md"
	if isFirst {
		label = "📁 Local AGENTS.md"
	}
	fmt.Printf("%s (%s):\n", label, relPath)
	fmt.Printf("  %s\n", strings.ReplaceAll(strings.TrimSpace(string(content)), "\n", "\n  "))
	return false
}

func printDomainLessons(rootDir, relPath string) {
	mappingPath := relPath
	if !strings.HasSuffix(mappingPath, "/") {
		mappingPath += "/"
	}
	domain := mapPathToDomain(mappingPath)

	if domain != "" {
		printSessionStrictLessons(rootDir, domain)
	}
}

func printSessionStrictLessons(rootDir string, domain string) {
	sections := collectLessonSections([]string{domain})
	if len(sections) == 0 {
		return
	}

	codingStandardsPath := filepath.Join(rootDir, ".agents/workflows/coding-standards.md")
	lessons, err := extractLessons(codingStandardsPath, sections)
	if err != nil {
		return
	}

	fmt.Printf("⚠️  Relevant Strict Lessons (%s):\n", strings.Join(sections, " & "))
	for _, section := range sections {
		if content, ok := lessons[section]; ok {
			fmt.Printf("  %s\n", strings.ReplaceAll(content, "\n", "\n  "))
		}
	}
}

func printRelatedADRs(rootDir, relPath string) {
	adrDir := filepath.Join(rootDir, "docs/adr")
	files, err := os.ReadDir(adrDir)
	if err != nil {
		return
	}

	pkgName := filepath.Base(relPath)
	var matches []string

	for _, f := range files {
		if isMatchableADR(f) {
			if match := scanADRForMatch(adrDir, f.Name(), relPath, pkgName); match != "" {
				matches = append(matches, match)
			}
		}
	}

	if len(matches) > 0 {
		fmt.Println("📚 Related ADRs:")
		for _, m := range matches {
			fmt.Printf("  %s\n", m)
		}
	}
}

func isMatchableADR(f os.DirEntry) bool {
	return !f.IsDir() && strings.HasSuffix(f.Name(), ".md") && f.Name() != "template.md"
}

func scanADRForMatch(adrDir, fileName, relPath, pkgName string) string {
	path := filepath.Join(adrDir, fileName)
	// #nosec G304 -- adr path is trusted
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	textContent := string(content)
	if strings.Contains(textContent, relPath) || strings.Contains(textContent, pkgName) {
		firstLine := ""
		scanner := bufio.NewScanner(strings.NewReader(textContent))
		if scanner.Scan() {
			firstLine = strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "# "))
		}
		return fmt.Sprintf("- %s: %s", fileName, firstLine)
	}
	return ""
}
