package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func SetPhaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set-phase <IDLE|RED|GREEN|REFACTOR>",
		Short: "Set the current workflow phase",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			phase := args[0]
			if phase != "IDLE" && phase != "RED" && phase != "GREEN" && phase != "REFACTOR" {
				fmt.Fprintf(os.Stderr, "Error: invalid phase '%s'\n", phase)
				os.Exit(1)
			}
			state := getState()
			state.Phase = phase
			saveState(state)
			if flagText {
				fmt.Printf("Phase set to: %s\n", phase)
			} else {
				printJSON(true, "set-phase", map[string]any{
					"message": fmt.Sprintf("Phase set to: %s", phase),
					"phase":   phase,
					"state":   state,
				}, nil)
			}
		},
	}
}
