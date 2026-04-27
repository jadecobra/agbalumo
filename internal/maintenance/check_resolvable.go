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

	skillDirs, err := scanSkillDirs(skillsDir)
	if err != nil {
		return violations
	}

	resolverSkills, err := parseResolverSkills(resolverPath)
	if err != nil {
		return violations
	}

	manifestSkills, _ := parseManifestSkills(manifestPath)
	return findCrossReferenceViolations(skillDirs, resolverSkills, manifestSkills)
}


func scanSkillDirs(dir string) (map[string]bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	skillDirs := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != ".DS_Store" {
			skillDirs[entry.Name()] = true
		}
	}
	return skillDirs, nil
}

func parseResolverSkills(path string) (map[string]bool, error) {
	// #nosec G304 -- maintenance utility reads trusted paths
	resolverBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	resolverContent := string(resolverBytes)

	procSkillsIdx := strings.Index(resolverContent, "## Procedural Skills")
	if procSkillsIdx == -1 {
		return nil, nil
	}

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

	return resolverSkills, nil
}

func parseManifestSkills(path string) (map[string]bool, error) {
	// #nosec G304 -- maintenance utility reads trusted paths
	manifestBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, err
	}

	manifestSkills := make(map[string]bool)
	for _, s := range manifest.Skills {
		manifestSkills[s.Name] = true
	}

	return manifestSkills, nil
}

func findCrossReferenceViolations(skillDirs, resolverSkills, manifestSkills map[string]bool) []string {
	var violations []string
	if resolverSkills != nil {
		violations = append(violations, findOrphanedSkills(skillDirs, resolverSkills)...)
		violations = append(violations, findDanglingEntries(skillDirs, resolverSkills)...)
	}
	if manifestSkills != nil {
		violations = append(violations, findUnregisteredSkills(skillDirs, manifestSkills)...)
	}
	return violations
}

func findOrphanedSkills(skillDirs, resolverSkills map[string]bool) []string {
	var violations []string
	for dir := range skillDirs {
		if !resolverSkills[dir] {
			violations = append(violations, "orphaned: "+dir)
		}
	}
	return violations
}

func findDanglingEntries(skillDirs, resolverSkills map[string]bool) []string {
	var violations []string
	for skill := range resolverSkills {
		if !skillDirs[skill] {
			violations = append(violations, "dangling: "+skill)
		}
	}
	return violations
}

func findUnregisteredSkills(skillDirs, manifestSkills map[string]bool) []string {
	var violations []string
	for dir := range skillDirs {
		if !manifestSkills[dir] {
			violations = append(violations, "unregistered: "+dir)
		}
	}
	return violations
}


