package agent_test

import (
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestCreateHandoff_YAMLFrontmatter(t *testing.T) {
	const testFile = "HANDOFF_test.md"

	state := &agent.State{
		Feature:      "test-feature",
		WorkflowType: "feature",
		Phase:        "RED",
	}
	state.Gates.RedTest = agent.GatePassed

	params := agent.HandoffParams{
		TargetPersona: "BackendEngineer",
		CurrentState:  state,
		Progress:      nil,
	}

	// Override the output path for test isolation
	err := agent.CreateHandoffToPath(params, testFile)
	if err != nil {
		t.Fatalf("CreateHandoffToPath returned error: %v", err)
	}
	defer func() { _ = agent.RemoveFile(testFile) }()

	content, readErr := agent.ReadFile(testFile)
	if readErr != nil {
		t.Fatalf("failed to read generated handoff: %v", readErr)
	}

	// Assert frontmatter block is present
	if !strings.HasPrefix(content, "---\n") {
		t.Errorf("expected HANDOFF.md to start with YAML frontmatter '---', got: %.50s", content)
	}

	// Assert required frontmatter fields exist
	for _, field := range []string{"target_persona:", "phase:", "feature:", "workflow_type:", "prior_gate_state:"} {
		if !strings.Contains(content, field) {
			t.Errorf("HANDOFF.md frontmatter missing required field: %s", field)
		}
	}

	// Assert the persona value is correct
	if !strings.Contains(content, "target_persona: BackendEngineer") {
		t.Errorf("expected target_persona to be 'BackendEngineer'")
	}

	// Assert human-readable body still present after frontmatter
	if !strings.Contains(content, "# 🤝 HANDOFF") {
		t.Errorf("expected HANDOFF.md body section to still be present")
	}
}
