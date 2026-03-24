package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print the current status of the harness",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			state := getState()
			if flagText {
				b, err := json.MarshalIndent(state, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling state: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(string(b))
			} else {
				printJSON(true, "status", map[string]any{"state": state}, nil)
			}
		},
	}
}
