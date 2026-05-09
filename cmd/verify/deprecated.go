package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

// Replaces Strict Lesson: Strict ViewModel Mandate
var deprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "Scan for deprecated patterns in Go files",
	Run: func(cmd *cobra.Command, args []string) {
		violations, err := maintenance.CheckDeprecatedPatterns(".")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if len(violations) > 0 {
			fmt.Printf("Found %d deprecated pattern violations:\n", len(violations))
			for _, v := range violations {
				fmt.Printf("%s:%d - %s: %s\n", v.File, v.Line, v.Pattern, v.Suggestion)
			}
			os.Exit(1)
		}

		fmt.Println("No deprecated patterns found.")
	},
}
