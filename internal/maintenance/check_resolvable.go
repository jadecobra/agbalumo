package maintenance

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Skills []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
	} `yaml:"skills"`
}

func CheckResolvable(skillsDir, resolverPath, manifestPath string) []string {
	var violations []string

	// 1. Scan skillsDir subdirectories
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return violations
	}

	skillDirs := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != ".DS_Store" {
			skillDirs[entry.Name()] = true
		}
	}

	// 2. Parse RESOLVER.md
	// #nosec G304 -- maintenance utility reads trusted paths
	resolverBytes, err := os.ReadFile(resolverPath)

	if err != nil {
		return violations
	}
	resolverContent := string(resolverBytes)

	// Extract Procedural Skills section
	procSkillsIdx := strings.Index(resolverContent, "## Procedural Skills")
	if procSkillsIdx != -1 {
		procSection := resolverContent[procSkillsIdx:]
		nextHeadingIdx := strings.Index(procSection[strings.Index(procSection, "\n")+1:], "##")
		if nextHeadingIdx != -1 {
			procSection = procSection[:strings.Index(procSection, "\n")+1+nextHeadingIdx]
		}

		// Parse rows via regex
		re := regexp.MustCompile(`\|\s*[^|]+\s*\|\s*([^|\s]+)\s*\|`)
		matches := re.FindAllStringSubmatch(procSection, -1)

		resolverSkills := make(map[string]bool)
		for _, match := range matches {
			if len(match) > 1 {
				path := match[1]
				// Extract skill name from path (e.g. .agents/skills/go-tdd/SKILL.md -> go-tdd)
				skillName := filepath.Base(filepath.Dir(path))
				if skillName != "." && skillName != "/" {
					resolverSkills[skillName] = true
				}
			}
		}

		// Cross-reference: Orphaned Skills (dir exists but no RESOLVER.md entry)
		for dir := range skillDirs {
			if !resolverSkills[dir] {
				violations = append(violations, "orphaned: "+dir)
			}
		}

		// Cross-reference: Dangling Entries (RESOLVER.md points to non-existent dir)
		for skill := range resolverSkills {
			if !skillDirs[skill] {
				violations = append(violations, "dangling: "+skill)
			}
		}
	}

	// 3. Parse verify-manifest.yaml
	// #nosec G304 -- maintenance utility reads trusted paths
	manifestBytes, err := os.ReadFile(manifestPath)

	if err == nil {
		var manifest Manifest
		if err := yaml.Unmarshal(manifestBytes, &manifest); err == nil {
			manifestSkills := make(map[string]bool)
			for _, s := range manifest.Skills {
				manifestSkills[s.Name] = true
			}

			// Cross-reference: Unregistered Skills (skill in skillsDir but not in manifest)
			for dir := range skillDirs {
				if !manifestSkills[dir] {
					violations = append(violations, "unregistered: "+dir)
				}
			}
		}
	}

	return violations
}
