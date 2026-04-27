package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillConformance_ValidSkill(t *testing.T) {
	tempDir := t.TempDir()
	skillDir := filepath.Join(tempDir, "go-tdd")
	if err := os.Mkdir(skillDir, 0700); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}

	skillFile := filepath.Join(skillDir, "SKILL.md")
	content := `---
name: Go TDD Workflow
description: Execute the RED-GREEN-REFACTOR cycle for Go projects
triggers:
  - "writing tests"
mutating: false
---
# Skill content`

	if err := os.WriteFile(skillFile, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	violations := SkillConformance(tempDir)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %v", len(violations), violations)
	}
}

func TestSkillConformance_MissingTriggers(t *testing.T) {
	tempDir := t.TempDir()
	skillDir := filepath.Join(tempDir, "go-tdd")
	if err := os.Mkdir(skillDir, 0700); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}

	skillFile := filepath.Join(skillDir, "SKILL.md")
	content := `---
name: Go TDD Workflow
description: Execute the RED-GREEN-REFACTOR cycle for Go projects
mutating: false
---
# Skill content`

	if err := os.WriteFile(skillFile, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	violations := SkillConformance(tempDir)
	found := false
	for _, v := range violations {
		if v == "go-tdd: missing 'triggers'" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation 'go-tdd: missing 'triggers'', got %v", violations)
	}
}

func TestSkillConformance_NoSkillFile(t *testing.T) {
	tempDir := t.TempDir()
	skillDir := filepath.Join(tempDir, "go-tdd")
	if err := os.Mkdir(skillDir, 0700); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}

	// No SKILL.md created in skillDir

	violations := SkillConformance(tempDir)
	found := false
	for _, v := range violations {
		if v == "go-tdd: missing SKILL.md" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation 'go-tdd: missing SKILL.md', got %v", violations)
	}
}

