package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var templateContractCmd = &cobra.Command{
	Use:   "template-contract",
	Short: "Enforces strictly typed ViewModel contracts for UI templates",
	Run: func(cmd *cobra.Command, args []string) {
		violations, err := maintenance.CheckTemplateKeyGaps(".")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if len(violations) > 0 {
			fmt.Printf("Found %d template contract violations:\n", len(violations))
			for _, v := range violations {
				fmt.Printf("%s:%d - %s\n  %s\n", v.File, v.Line, v.Reason, v.Content)
			}
			os.Exit(1)
		}

		fmt.Println("Template contract check passed.")
	},
}
