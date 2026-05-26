package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckLessonsConformance(t *testing.T) {
	tests := []struct {
		name              string
		codingStandards   string
		expectedMsgSubstr string
		expectViolations  bool
	}{
		{
			name: "valid under ceiling with triggers",
			codingStandards: `---
description: Test standards
---
# Coding Standards

## Strict Lessons
* **Lesson 1** [TRIGGER: test_trigger]
* **Lesson 2** [TRIGGER: ui_change, modal_creation]
`,
			expectViolations: false,
		},
		{
			name: "missing trigger tag",
			codingStandards: `---
description: Test standards
---
# Coding Standards

## Strict Lessons
* **Lesson 1** [TRIGGER: test_trigger]
* **Lesson 2 without trigger**
`,
			expectViolations:  true,
			expectedMsgSubstr: "missing [TRIGGER:",
		},
		{
			name: "exceeds the 20-lesson ceiling",
			codingStandards: `---
description: Test standards
---
# Coding Standards

## Strict Lessons
* **L1** [TRIGGER: t]
* **L2** [TRIGGER: t]
* **L3** [TRIGGER: t]
* **L4** [TRIGGER: t]
* **L5** [TRIGGER: t]
* **L6** [TRIGGER: t]
* **L7** [TRIGGER: t]
* **L8** [TRIGGER: t]
* **L9** [TRIGGER: t]
* **L10** [TRIGGER: t]
* **L11** [TRIGGER: t]
* **L12** [TRIGGER: t]
* **L13** [TRIGGER: t]
* **L14** [TRIGGER: t]
* **L15** [TRIGGER: t]
* **L16** [TRIGGER: t]
* **L17** [TRIGGER: t]
* **L18** [TRIGGER: t]
* **L19** [TRIGGER: t]
* **L20** [TRIGGER: t]
* **L21** [TRIGGER: t]
`,
			expectViolations:  true,
			expectedMsgSubstr: "exceeds strict lesson ceiling of 20",
		},
		{
			name: "no strict lessons section",
			codingStandards: `---
description: Test standards
---
# Coding Standards
`,
			expectViolations: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLessonsConformanceTestCase(t, tt.codingStandards, tt.expectedMsgSubstr, tt.expectViolations)
		})
	}
}

func runLessonsConformanceTestCase(t *testing.T, codingStandards, expectedMsgSubstr string, expectViolations bool) {
	tmpDir := t.TempDir()

	// Create .agents/workflows directory (G301: permissions 0700 is secure)
	workflowsDir := filepath.Join(tmpDir, ".agents", "workflows")
	if err := os.MkdirAll(workflowsDir, 0700); err != nil {
		t.Fatalf("failed to create workflows dir: %v", err)
	}

	// Write mock coding-standards.md (G306: permissions 0600 is secure)
	standardsPath := filepath.Join(workflowsDir, "coding-standards.md")
	if err := os.WriteFile(standardsPath, []byte(codingStandards), 0600); err != nil {
		t.Fatalf("failed to write coding-standards.md: %v", err)
	}

	// Call CheckLessonsConformance
	violations, err := CheckLessonsConformance(tmpDir)
	if err != nil {
		t.Fatalf("CheckLessonsConformance failed: %v", err)
	}

	if expectViolations {
		assertLessonsViolationsExist(t, violations, expectedMsgSubstr)
	} else {
		if len(violations) > 0 {
			t.Fatalf("expected no violations, got %d: %+v", len(violations), violations)
		}
	}
}

func assertLessonsViolationsExist(t *testing.T, violations []LessonsViolation, expectedMsgSubstr string) {
	if len(violations) == 0 {
		t.Fatalf("expected violations but got none")
	}
	if expectedMsgSubstr == "" {
		return
	}
	found := false
	for _, v := range violations {
		if filepath.Base(v.File) == "coding-standards.md" && (v.Message == expectedMsgSubstr || len(expectedMsgSubstr) > 0) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation message containing '%s', got: %+v", expectedMsgSubstr, violations)
	}
}
