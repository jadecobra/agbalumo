package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SetPhaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set-phase <IDLE|RED|GREEN|REFACTOR>",
		Short: "Set the current workflow phase",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			phase := args[0]
			if phase != "IDLE" && phase != "RED" && phase != "GREEN" && phase != "REFACTOR" {
				return fmt.Errorf("invalid phase '%s'", phase)
			}
			state, err := getState()
			if err != nil {
				return err
			}
			state.Phase = phase
			if err := saveState(state); err != nil {
				return err
			}
			if flagText {
				fmt.Printf("Phase set to: %s\n", phase)
			} else {
				printJSON(true, "set-phase", map[string]any{
					"message": fmt.Sprintf("Phase set to: %s", phase),
					"phase":   phase,
					"state":   state,
				}, nil)
			}
			return nil
		},
	}
}
