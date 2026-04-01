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
	var sb strings.Builder

	// Prepend YAML frontmatter
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "target_persona: %s\n", params.TargetPersona)
	fmt.Fprintf(&sb, "feature: %s\n", params.CurrentState.Feature)
	fmt.Fprintf(&sb, "workflow_type: %s\n", params.CurrentState.WorkflowType)
	fmt.Fprintf(&sb, "phase: %s\n", params.CurrentState.Phase)

	// Capture relevant gate states for validation
	fmt.Fprintf(&sb, "prior_gate_state:\n")
	fmt.Fprintf(&sb, "  red_test: %s\n", params.CurrentState.Gates.RedTest)
	fmt.Fprintf(&sb, "  implementation: %s\n", params.CurrentState.Gates.Implementation)
	fmt.Fprintf(&sb, "  lint: %s\n", params.CurrentState.Gates.Lint)
	fmt.Fprintf(&sb, "---\n\n")

	fmt.Fprintf(&sb, "# 🤝 HANDOFF: Transition to %s\n\n", params.TargetPersona)
	fmt.Fprintf(&sb, "> [!IMPORTANT]\n")
	fmt.Fprintf(&sb, "> This is a context-bridge for the **%s** persona. \n", params.TargetPersona)
	fmt.Fprintf(&sb, "> To proceed: Open a NEW chat window, tag `@%s`, and run `/resume`.\n\n", params.TargetPersona)

	fmt.Fprintf(&sb, "## 📋 Project Context\n")
	fmt.Fprintf(&sb, "- **Feature**: %s\n", params.CurrentState.Feature)
	fmt.Fprintf(&sb, "- **Workflow**: %s\n", params.CurrentState.WorkflowType)
	fmt.Fprintf(&sb, "- **Current Phase**: %s\n\n", params.CurrentState.Phase)

	fmt.Fprintf(&sb, "## 🚀 Outgoing Persona Summary\n")
	fmt.Fprintf(&sb, "The previous task has reached a milestone. \n")
	if params.Progress != nil && len(params.Progress.Features) > 0 {
		latest := params.Progress.Features[len(params.Progress.Features)-1]
		fmt.Fprintf(&sb, "Latest Completed Category: **%s**\n", latest.Category)
	}
	fmt.Fprintf(&sb, "\n")

	fmt.Fprintf(&sb, "## 🛠️ Ground Truth Files\n")
	fmt.Fprintf(&sb, "The following files contain the canonical state for this feature. Review them immediately after resuming:\n")
	wd, _ := util.SafeGetwd()
	fmt.Fprintf(&sb, "- [implementation_plan.md](file:///%s/implementation_plan.md)\n", wd)
	fmt.Fprintf(&sb, "- [progress.md](file:///%s/.tester/tasks/progress.md)\n", wd)
	fmt.Fprintf(&sb, "- [state.json](file:///%s/.agents/state.json)\n\n", wd)

	fmt.Fprintf(&sb, "## 🎯 Next Objective\n")
	fmt.Fprintf(&sb, "As the **%s**, your goal is to transition the project to the next state as defined in the `implementation_plan.md`.\n", params.TargetPersona)
	fmt.Fprintf(&sb, "1. Run `task lint` and `task test` to verify the current baseline.\n")
	fmt.Fprintf(&sb, "2. Follow the `/build-feature` pipeline for the next phase.\n")

	return util.SafeWriteFile(path, []byte(sb.String()))
}
