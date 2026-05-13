package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Replaces Strict Lesson: UI Regression Verification
var browserCmd = &cobra.Command{
	Use:                "browser",
	Short:              "Execute Playwright end-to-end UI verification tests",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎭 Running Playwright E2E tests...")
		if len(args) > 0 {
			npmArgs := append([]string{"run", "test:e2e", "--"}, args...)
			return runCmd("npm", npmArgs...)
		}
		return runCmd("npm", "run", "test:e2e")
	},
}
