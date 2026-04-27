package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckResolvable_AllResolved(t *testing.T) {
	tempDir := t.TempDir()
	skillsDir := filepath.Join(tempDir, "skills")
	if err := os.Mkdir(skillsDir, 0700); err != nil {
		t.Fatalf("failed to create skills dir: %v", err)
	}

	// Create skill dir
	goTddDir := filepath.Join(skillsDir, "go-tdd")
	if err := os.Mkdir(goTddDir, 0700); err != nil {
		t.Fatalf("failed to create go-tdd dir: %v", err)
	}

	// Create RESOLVER.md
	resolverPath := filepath.Join(tempDir, "RESOLVER.md")
	resolverContent := `
## Procedural Skills
| Trigger | Skill |
|---------|-------|
| Writing tests | .agents/skills/go-tdd/SKILL.md |
`
	if err := os.WriteFile(resolverPath, []byte(resolverContent), 0600); err != nil {
		t.Fatalf("failed to write RESOLVER.md: %v", err)
	}

	// Create verify-manifest.yaml
	manifestPath := filepath.Join(tempDir, "verify-manifest.yaml")
	manifestContent := `
skills:
  - name: go-tdd
    path: .agents/skills/go-tdd/SKILL.md
`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0600); err != nil {
		t.Fatalf("failed to write verify-manifest.yaml: %v", err)
	}

	violations := CheckResolvable(skillsDir, resolverPath, manifestPath)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %v", len(violations), violations)
	}
}

func TestCheckResolvable_OrphanedSkill(t *testing.T) {
	tempDir := t.TempDir()
	skillsDir := filepath.Join(tempDir, "skills")
	if err := os.Mkdir(skillsDir, 0700); err != nil {
		t.Fatalf("failed to create skills dir: %v", err)
	}

	// Create skill dir
	goTddDir := filepath.Join(skillsDir, "go-tdd")
	if err := os.Mkdir(goTddDir, 0700); err != nil {
		t.Fatalf("failed to create go-tdd dir: %v", err)
	}

	// Create RESOLVER.md (empty/no go-tdd)
	resolverPath := filepath.Join(tempDir, "RESOLVER.md")
	resolverContent := `
## Procedural Skills
| Trigger | Skill |
|---------|-------|
`
	if err := os.WriteFile(resolverPath, []byte(resolverContent), 0600); err != nil {
		t.Fatalf("failed to write RESOLVER.md: %v", err)
	}

	// Create verify-manifest.yaml
	manifestPath := filepath.Join(tempDir, "verify-manifest.yaml")
	manifestContent := `
skills:
  - name: go-tdd
    path: .agents/skills/go-tdd/SKILL.md
`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0600); err != nil {
		t.Fatalf("failed to write verify-manifest.yaml: %v", err)
	}

	violations := CheckResolvable(skillsDir, resolverPath, manifestPath)
	found := false
	for _, v := range violations {
		if v == "orphaned: go-tdd" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation 'orphaned: go-tdd', got %v", violations)
	}
}

func TestCheckResolvable_DanglingEntry(t *testing.T) {
	tempDir := t.TempDir()
	skillsDir := filepath.Join(tempDir, "skills")
	if err := os.Mkdir(skillsDir, 0700); err != nil {
		t.Fatalf("failed to create skills dir: %v", err)
	}

	// Create RESOLVER.md with dangling entry
	resolverPath := filepath.Join(tempDir, "RESOLVER.md")
	resolverContent := `
## Procedural Skills
| Trigger | Skill |
|---------|-------|
| Writing tests | .agents/skills/foo/SKILL.md |
`
	if err := os.WriteFile(resolverPath, []byte(resolverContent), 0600); err != nil {
		t.Fatalf("failed to write RESOLVER.md: %v", err)
	}

	// Create verify-manifest.yaml
	manifestPath := filepath.Join(tempDir, "verify-manifest.yaml")
	manifestContent := `
skills:
  - name: foo
    path: .agents/skills/foo/SKILL.md
`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0600); err != nil {
		t.Fatalf("failed to write verify-manifest.yaml: %v", err)
	}

	violations := CheckResolvable(skillsDir, resolverPath, manifestPath)
	found := false
	for _, v := range violations {
		if v == "dangling: foo" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation 'dangling: foo', got %v", violations)
	}
}

func TestCheckResolvable_MissingManifest(t *testing.T) {
	tempDir := t.TempDir()
	skillsDir := filepath.Join(tempDir, "skills")
	if err := os.Mkdir(skillsDir, 0700); err != nil {
		t.Fatalf("failed to create skills dir: %v", err)
	}

	// Create skill dir
	goTddDir := filepath.Join(skillsDir, "go-tdd")
	if err := os.Mkdir(goTddDir, 0700); err != nil {
		t.Fatalf("failed to create go-tdd dir: %v", err)
	}

	// Create RESOLVER.md
	resolverPath := filepath.Join(tempDir, "RESOLVER.md")
	resolverContent := `
## Procedural Skills
| Trigger | Skill |
|---------|-------|
| Writing tests | .agents/skills/go-tdd/SKILL.md |
`
	if err := os.WriteFile(resolverPath, []byte(resolverContent), 0600); err != nil {
		t.Fatalf("failed to write RESOLVER.md: %v", err)
	}

	// Create verify-manifest.yaml (empty skills)
	manifestPath := filepath.Join(tempDir, "verify-manifest.yaml")
	manifestContent := `
skills:
`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0600); err != nil {
		t.Fatalf("failed to write verify-manifest.yaml: %v", err)
	}

	violations := CheckResolvable(skillsDir, resolverPath, manifestPath)
	found := false
	for _, v := range violations {
		if v == "unregistered: go-tdd" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation 'unregistered: go-tdd', got %v", violations)
	}
}
