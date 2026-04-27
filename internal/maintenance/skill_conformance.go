package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SkillConformance checks that each skill directory in skillsDir contains a valid SKILL.md.
func SkillConformance(skillsDir string) []string {
	var violations []string
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return []string{fmt.Sprintf("failed to read skills dir %s: %v", skillsDir, err)}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillName := entry.Name()
		skillDir := filepath.Join(skillsDir, skillName)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			violations = append(violations, fmt.Sprintf("%s: missing SKILL.md", skillName))
			continue
		}

		// Read and parse SKILL.md
		// #nosec G304 -- path is derived from trusted repo structure
		content, err := os.ReadFile(skillFile)
		if err != nil {
			violations = append(violations, fmt.Sprintf("%s: failed to read SKILL.md: %v", skillName, err))
			continue
		}

		yamlStr := extractFrontmatter(string(content))

		var fm struct {
			Name        *string   `yaml:"name"`
			Description *string   `yaml:"description"`
			Triggers    *[]string `yaml:"triggers"`
			Mutating    *bool     `yaml:"mutating"`
		}

		if err := yaml.Unmarshal([]byte(yamlStr), &fm); err != nil {
			violations = append(violations, fmt.Sprintf("%s: invalid YAML frontmatter: %v", skillName, err))
			continue
		}

		if fm.Name == nil || *fm.Name == "" {
			violations = append(violations, fmt.Sprintf("%s: missing 'name'", skillName))
		}
		if fm.Description == nil || *fm.Description == "" {
			violations = append(violations, fmt.Sprintf("%s: missing 'description'", skillName))
		}
		if fm.Triggers == nil {
			violations = append(violations, fmt.Sprintf("%s: missing 'triggers'", skillName))
		}
		if fm.Mutating == nil {
			violations = append(violations, fmt.Sprintf("%s: missing 'mutating'", skillName))
		}
	}

	return violations
}

func extractFrontmatter(content string) string {
	lines := strings.Split(content, "\n")
	var yamlLines []string
	inYAML := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inYAML {
				inYAML = true
				continue
			} else {
				break
			}
		}
		if inYAML {
			yamlLines = append(yamlLines, line)
		}
	}
	return strings.Join(yamlLines, "\n")
}
