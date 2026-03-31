package commands

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"github.com/spf13/cobra"
)

func HandoffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "handoff [target_persona]",
		Short: "Generate a HANDOFF.md context-bridge for a new chat window",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := getState()
			if err != nil {
				return err
			}

			targetPersona := "BackendEngineer" // Default
			if len(args) > 0 {
				targetPersona = args[0]
			} else {
				// Intelligent mapping based on current phase (expandable)
				switch state.Phase {
				case "RED":
					targetPersona = "BackendEngineer"
				case "GREEN":
					targetPersona = "ChiefCritic"
				case "REFACTOR":
					targetPersona = "SecurityEngineer"
				}
			}

			// Load progress
			progressPath := ".tester/tasks/progress.md"
			progressData, _ := util.SafeReadFile(progressPath)
			tracker, _ := agent.ParseMarkdownTracker(string(progressData))

			params := agent.HandoffParams{
				TargetPersona: targetPersona,
				CurrentState:  state,
				Progress:      &tracker,
			}

			if err := agent.CreateHandoff(params); err != nil {
				return fmt.Errorf("failed to create handoff: %w", err)
			}

			if flagText {
				fmt.Printf("🤝 Handoff bridge created: %s\n", agent.HandoffFile)
				fmt.Printf("➡️  Target Persona: %s\n", targetPersona)
				fmt.Printf("👉 ACTION: Open a NEW conversation, tag @%s, and run `/resume`.\n", targetPersona)
			} else {
				printJSON(true, "handoff", map[string]any{
					"handoff_file":   agent.HandoffFile,
					"target_persona": targetPersona,
					"message":        fmt.Sprintf("Handoff bridge created for %s", targetPersona),
				}, nil)
			}

			return nil
		},
	}
}
