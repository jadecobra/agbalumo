package maintenance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type VerifyManifest struct {
	Commands []ManifestCommand `yaml:"commands"`
	Skills   []ManifestSkill   `yaml:"skills"`
}

type ManifestCommand struct {
	Name        string `yaml:"name"`
	Trigger     string `yaml:"trigger"`
	Description string `yaml:"description"`
}

type ManifestSkill struct {
	Name    string `yaml:"name"`
	Trigger string `yaml:"trigger"`
	Path    string `yaml:"path"`
}

const (
	domainHandler    = "handler"
	domainDomain     = "domain"
	domainRepository = "repository"
	domainService    = "service"
	domainUI         = "ui"
	domainCI         = "ci"
	domainTesting    = "testing"
)

// RunPreflight dumps active rules relevant to staged/modified files.
func RunPreflight(rootDir string) error {
	modifiedFiles, err := getGitModifiedFiles(rootDir)
	if err != nil {
		return err
	}

	domains := collectDomains(modifiedFiles)
	if len(domains) == 0 {
		fmt.Println("No modified domains detected.")
		return nil
	}

	fmt.Println("📋 Preflight Context for this session")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Modified domains: [%s]\n", strings.Join(domains, ", "))

	stagedOut, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err == nil {
		stagedFiles := strings.Split(string(stagedOut), "\n")
		hasUI := false
		hasCI := false
		hasTest := false
		hasEnv := false

		for _, f := range stagedFiles {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			if strings.HasSuffix(f, ".js") || strings.HasSuffix(f, ".css") || strings.HasSuffix(f, ".html") || strings.HasSuffix(f, ".tmpl") || strings.HasPrefix(f, "ui/") {
				hasUI = true
			}
			if strings.HasPrefix(f, ".github/") || strings.Contains(f, "workflow") || strings.HasPrefix(f, "scripts/") || strings.HasPrefix(f, "cmd/verify/") || strings.Contains(filepath.Base(f), "ci") {
				hasCI = true
			}
			if strings.HasSuffix(f, "_test.go") || strings.Contains(f, "test") {
				hasTest = true
			}
			if strings.HasPrefix(filepath.Base(f), ".env") || strings.Contains(f, "secret") || strings.Contains(f, "security") {
				hasEnv = true
			}
		}

		if hasUI {
			fmt.Println("Read: coding-standards.md → ### UI & Frontend")
		}
		if hasCI {
			fmt.Println("Read: coding-standards.md → ### CI & Infrastructure")
		}
		if hasTest {
			fmt.Println("Read: coding-standards.md → ### Testing")
		}
		if hasEnv {
			fmt.Println("Read: coding-standards.md → ### Security & Environment")
		}
	}

	printPackageConstraints(rootDir, domains)
	printStrictLessons(rootDir, domains)
	printSkillsAndCommands(rootDir, modifiedFiles)
	printInvariants(rootDir)

	return nil
}

func printPackageConstraints(rootDir string, domains []string) {
	for _, domain := range domains {
		var agentsPaths []string
		switch domain {
		case domainRepository:
			agentsPaths = []string{"internal/repository/AGENTS.md"}
		case domainHandler:
			agentsPaths = []string{"internal/handler/AGENTS.md", "internal/module/AGENTS.md"}
		case domainService:
			agentsPaths = []string{"internal/service/AGENTS.md"}
		case domainUI:
			agentsPaths = []string{"internal/ui/AGENTS.md", "ui/AGENTS.md"}
		case domainCI:
			agentsPaths = []string{".github/AGENTS.md"}
		case domainDomain:
			agentsPaths = []string{"internal/domain/AGENTS.md"}
		}

		for _, agentsPath := range agentsPaths {
			fullPath := filepath.Join(rootDir, agentsPath)
			// #nosec G304 -- path is derived from trusted repo structure
			if content, err := os.ReadFile(fullPath); err == nil {
				fmt.Printf("📁 Package Constraints (%s):\n", agentsPath)
				fmt.Printf("  %s\n", strings.ReplaceAll(strings.TrimSpace(string(content)), "\n", "\n  "))
			}
		}
	}
}

func printStrictLessons(rootDir string, domains []string) {
	sections := collectLessonSections(domains)
	if len(sections) == 0 {
		return
	}

	codingStandardsPath := filepath.Join(rootDir, ".agents/workflows/coding-standards.md")
	lessons, err := extractLessons(codingStandardsPath, sections)
	if err != nil {
		return
	}

	for _, section := range sections {
		if content, ok := lessons[section]; ok {
			fmt.Printf("⚠️  Active Strict Lessons (%s):\n", section)
			fmt.Printf("  %s\n", strings.ReplaceAll(content, "\n", "\n  "))
		}
	}
}

func printInvariants(rootDir string) {
	invariantsPath := filepath.Join(rootDir, ".agents/invariants.json")
	// #nosec G304 -- invariants path is trusted project metadata
	content, err := os.ReadFile(invariantsPath)
	if err != nil {
		return
	}

	fmt.Println("🔧 Invariants:")
	var inv map[string]interface{}
	if err := json.Unmarshal(content, &inv); err == nil {
		keys := make([]string, 0, len(inv))
		for k := range inv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %s: %v\n", k, inv[k])
		}
	}
}

func printSkillsAndCommands(rootDir string, modifiedFiles []string) {
	manifest, err := loadVerifyManifest(rootDir)
	if err != nil {
		return
	}

	matchedTriggers := getMatchedTriggers(modifiedFiles)
	printMatchedSkills(manifest.Skills, matchedTriggers)
	printMatchedCommands(manifest.Commands, matchedTriggers)
}

func loadVerifyManifest(rootDir string) (*VerifyManifest, error) {
	manifestPath := filepath.Join(rootDir, ".agents/verify-manifest.yaml")
	// #nosec G304 -- manifest path is trusted project metadata
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest VerifyManifest
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func printMatchedSkills(skills []ManifestSkill, matchedTriggers map[string]bool) {
	var results []string
	for _, s := range skills {
		if isManifestItemMatched(s.Trigger, matchedTriggers) {
			results = append(results, fmt.Sprintf("- %s → %s", s.Name, s.Path))
		}
	}

	if len(results) > 0 {
		fmt.Println("📖 Relevant Skills:")
		for _, r := range results {
			fmt.Printf("  %s\n", r)
		}
	}
}

func printMatchedCommands(commands []ManifestCommand, matchedTriggers map[string]bool) {
	var results []string
	for _, c := range commands {
		if isManifestItemMatched(c.Trigger, matchedTriggers) {
			desc := c.Description
			if desc == "" {
				desc = "No description"
			}
			results = append(results, fmt.Sprintf("- %s (%s)", c.Name, desc))
		}
	}

	if len(results) > 0 {
		fmt.Println("🛠️  Relevant Verify Commands:")
		for _, r := range results {
			fmt.Printf("  %s\n", r)
		}
	}
}

func isManifestItemMatched(triggerStr string, matchedTriggers map[string]bool) bool {
	triggers := strings.Split(triggerStr, ",")
	for _, t := range triggers {
		if matchedTriggers[strings.TrimSpace(t)] {
			return true
		}
	}
	return false
}

func getMatchedTriggers(modifiedFiles []string) map[string]bool {
	matched := make(map[string]bool)
	for _, f := range modifiedFiles {
		if strings.HasSuffix(f, "_test.go") || strings.HasPrefix(f, "internal/") {
			matched["test_authoring"] = true
			matched["feature_implementation"] = true
			matched["bug_fix"] = true
		}
		if isUIDomain(f) {
			matched["ui_change"] = true
			matched["browser_subagent"] = true
		}
	}
	return matched
}

func getGitModifiedFiles(rootDir string) ([]string, error) {
	cmd1 := exec.Command("git", "diff", "--name-only", "HEAD")
	cmd1.Dir = rootDir
	out1, _ := cmd1.Output()

	cmd2 := exec.Command("git", "diff", "--cached", "--name-only")
	cmd2.Dir = rootDir
	out2, _ := cmd2.Output()

	files := make(map[string]bool)
	for _, line := range strings.Split(string(out1)+string(out2), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files[line] = true
		}
	}

	var result []string
	for f := range files {
		result = append(result, f)
	}
	sort.Strings(result)
	return result, nil
}

func collectDomains(files []string) []string {
	domainMap := make(map[string]bool)
	for _, f := range files {
		d := mapPathToDomain(f)
		if d != "" {
			domainMap[d] = true
		}
	}

	var domains []string
	for d := range domainMap {
		domains = append(domains, d)
	}
	sort.Strings(domains)
	return domains
}

func mapPathToDomain(path string) string {
	switch {
	case strings.HasSuffix(path, "_test.go"):
		return domainTesting
	case strings.HasPrefix(path, "internal/repository/") || strings.HasSuffix(path, ".sql"):
		return domainRepository
	case strings.HasPrefix(path, "internal/handler/") || strings.HasPrefix(path, "internal/module/"):
		return domainHandler
	case strings.HasPrefix(path, "internal/service/"):
		return domainService
	case isUIDomain(path):
		return domainUI
	case isCIDomain(path):
		return domainCI
	case strings.HasPrefix(path, "internal/domain/"):
		return domainDomain
	default:
		return ""
	}
}

func isUIDomain(path string) bool {
	return strings.HasPrefix(path, "ui/") ||
		strings.HasSuffix(path, ".js") ||
		strings.HasSuffix(path, ".css") ||
		strings.HasSuffix(path, ".html")
}

func isCIDomain(path string) bool {
	return strings.HasPrefix(path, ".github/") ||
		strings.HasPrefix(path, "scripts/") ||
		strings.HasPrefix(path, "cmd/verify/")
}

func collectLessonSections(domains []string) []string {
	sectionMap := make(map[string]bool)
	for _, d := range domains {
		switch d {
		case domainHandler, domainUI:
			sectionMap["UI & Frontend"] = true
		case domainCI:
			sectionMap["CI & Infrastructure"] = true
		case domainTesting:
			sectionMap["Testing"] = true
		}
	}

	var sections []string
	for s := range sectionMap {
		sections = append(sections, s)
	}
	sort.Strings(sections)
	return sections
}

func extractLessons(path string, sections []string) (map[string]string, error) {
	// #nosec G304 -- path is derived from trusted repo structure
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)

	targetSections := make(map[string]bool)
	for _, s := range sections {
		targetSections["### "+s] = true
	}

	var currentSection string
	var currentContent strings.Builder
	recording := false

	for scanner.Scan() {
		line := scanner.Text()
		recording, currentSection = processLessonLine(line, recording, currentSection, targetSections, result, &currentContent)
	}
	if recording {
		result[currentSection] = strings.TrimSpace(currentContent.String())
	}

	return result, nil
}

func processLessonLine(line string, recording bool, currentSection string, targetSections map[string]bool, result map[string]string, content *strings.Builder) (bool, string) {
	if strings.HasPrefix(line, "### ") {
		if recording {
			result[currentSection] = strings.TrimSpace(content.String())
			content.Reset()
		}
		if targetSections[line] {
			return true, strings.TrimPrefix(line, "### ")
		}
		return false, ""
	}

	if recording {
		if strings.HasPrefix(line, "## ") {
			result[currentSection] = strings.TrimSpace(content.String())
			content.Reset()
			return false, ""
		}
		content.WriteString(line + "\n")
	}
	return recording, currentSection
}
