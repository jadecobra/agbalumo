package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func HelloCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "hello",
		Short:  "Protocol test command for agent handoff verification",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, Agent! Protocol test initialized.")
		},
	}
}
