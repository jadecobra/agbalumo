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
				b, err := json.MarshalIndent(state, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshaling state: %w", err)
				}
				fmt.Println(string(b))
			} else {
				printJSON(true, "status", map[string]any{"state": state}, nil)
			}
			return nil
		},
	}
}
