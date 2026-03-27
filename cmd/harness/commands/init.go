package commands

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init <feature_name> [workflow_type]",
		Short: "Initialize the harness tracking structure",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			feature := args[0]
			workflowType := agent.WorkflowFeature
			if len(args) > 1 {
				workflowType = args[1]
			}
			if workflowType != agent.WorkflowFeature && workflowType != agent.WorkflowBugfix && workflowType != agent.WorkflowRefactor {
				return fmt.Errorf("invalid workflow type '%s'", workflowType)
			}

			state := &agent.State{
				Feature:      feature,
				WorkflowType: workflowType,
				Phase:        "IDLE",
				Gates: agent.Gates{
					RedTest:             agent.GatePending,
					ApiSpec:             agent.GatePending,
					Implementation:      agent.GatePending,
					Lint:                agent.GatePending,
					Coverage:            agent.GatePending,
					BrowserVerification: agent.GatePending,
					TemplateDrift:       agent.GatePending,
				},
			}
			if err := saveState(state); err != nil {
				return err
			}
			if flagText {
				fmt.Printf("Workflow initialized for %s: %s\n", workflowType, feature)
				if err := summarizeProgress(); err != nil {
					return fmt.Errorf("failed to summarize progress: %w", err)
				}
			} else {
				printJSON(true, "init", map[string]any{
					"message":      fmt.Sprintf("Workflow initialized for %s: %s", workflowType, feature),
					"feature":      feature,
					"workflowType": workflowType,
					"state":        state,
				}, nil)
			}
			return nil
		},
	}
}
