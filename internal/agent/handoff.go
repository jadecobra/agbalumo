package agent

import (
	"fmt"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

const HandoffFile = ".tester/tasks/HANDOFF.md"

// HandoffParams contains the data required to generate a handoff file.
type HandoffParams struct {
	TargetPersona string
	CurrentState  *State
	Progress      *ProgressTracker
}

// ReadFile is a convenience wrapper for util.SafeReadFile.
func ReadFile(filename string) (string, error) {
	b, err := util.SafeReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// RemoveFile is a convenience wrapper for util.SafeRemove.
func RemoveFile(filename string) error {
	return util.SafeRemove(filename)
}

// CreateHandoff creates the HANDOFF.md file.
func CreateHandoff(params HandoffParams) error {
	return CreateHandoffToPath(params, HandoffFile)
}

// CreateHandoffToPath creates a handoff file at the specified path with YAML frontmatter.
func CreateHandoffToPath(params HandoffParams, path string) error {
	if params.CurrentState == nil {
		return fmt.Errorf("CreateHandoffToPath: CurrentState is required")
	}

	var sb strings.Builder

	// YAML Frontmatter
	sb.WriteString("---\n")
	fmt.Fprintf(&sb, "target_persona: %s\n", params.TargetPersona)
	fmt.Fprintf(&sb, "phase: %s\n", params.CurrentState.Phase)
	fmt.Fprintf(&sb, "feature: %s\n", params.CurrentState.Feature)
	fmt.Fprintf(&sb, "workflow_type: %s\n", params.CurrentState.WorkflowType)
	fmt.Fprintf(&sb, "prior_gate_state: red-test=%s api-spec=%s implementation=%s\n",
		string(params.CurrentState.Gates.RedTest),
		string(params.CurrentState.Gates.ApiSpec),
		string(params.CurrentState.Gates.Implementation))
	sb.WriteString("---\n\n")

	// Markdown Body
	fmt.Fprintf(&sb, "# 🤝 HANDOFF: Transition to %s\n\n", params.TargetPersona)
	fmt.Fprintf(&sb, "> [!IMPORTANT]\n")
	fmt.Fprintf(&sb, "> This is a context-bridge for the **%s** persona.\n", params.TargetPersona)
	fmt.Fprintf(&sb, "> To proceed: Open a NEW chat window and run `/resume`.\n\n")

	sb.WriteString("## 📋 Project Context\n")
	fmt.Fprintf(&sb, "- **Feature**: %s\n", params.CurrentState.Feature)
	fmt.Fprintf(&sb, "- **Workflow**: %s\n", params.CurrentState.WorkflowType)
	fmt.Fprintf(&sb, "- **Current Phase**: %s\n\n", params.CurrentState.Phase)


	sb.WriteString("## 🎯 Next Objective\n")
	fmt.Fprintf(&sb, "As the **%s**, your goal is to transition the project to the next state.\n", params.TargetPersona)
	sb.WriteString("1. Run `task lint` and `task test` to verify the current baseline.\n")
	sb.WriteString("2. Follow the `/build-feature` pipeline for the next phase.\n")

	return util.SafeWriteFile(path, []byte(sb.String()))
}






