package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print the current status of the harness",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := getState()
			if err != nil {
				return err
			}
			if flagText {
				fmt.Printf("--- Current Workflow Status ---\n")
				fmt.Printf("Feature: %s\n", state.Feature)
				fmt.Printf("Phase: %s\n", state.Phase)
				fmt.Printf("Workflow: %s\n", state.WorkflowType)
				fmt.Printf("Updated: %s\n", state.UpdatedAt)
				fmt.Printf("-------------------------------\n")
				b, _ := json.MarshalIndent(state.Gates, "", "  ")
				fmt.Printf("Gates:\n%s\n", string(b))
			} else {
				printJSON(true, "status", map[string]any{"state": state}, nil)
			}
			return nil
		},
	}
}
