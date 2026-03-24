package commands

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init <feature_name> [workflow_type]",
		Short: "Initialize the harness tracking structure",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			feature := args[0]
			workflowType := "feature"
			if len(args) > 1 {
				workflowType = args[1]
			}
			if workflowType != "feature" && workflowType != "bugfix" && workflowType != "refactor" {
				fmt.Fprintf(os.Stderr, "Error: invalid workflow type '%s'\n", workflowType)
				os.Exit(1)
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
				},
			}
			saveState(state)
			if flagText {
				fmt.Printf("Workflow initialized for %s: %s\n", workflowType, feature)
				summarizeProgress()
			} else {
				printJSON(true, "init", map[string]any{
					"message":      fmt.Sprintf("Workflow initialized for %s: %s", workflowType, feature),
					"feature":      feature,
					"workflowType": workflowType,
					"state":        state,
				}, nil)
			}
		},
	}
}
