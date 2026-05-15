package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var sandboxParityCmd = &cobra.Command{
	Use:   "sandbox-parity",
	Short: "Verify all ui_components.html partials are documented in sandbox.html",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking Sandbox Parity...")

		violations, err := maintenance.CheckSandboxParity(".")
		if err != nil {
			return fmt.Errorf("failed to check sandbox parity: %w", err)
		}

		if len(violations) == 0 {
			fmt.Println("✅ All components in ui_components.html are documented in sandbox.html.")
			return nil
		}

		fmt.Printf("\n❌ Found %d undocumented components:\n", len(violations))
		for _, v := range violations {
			fmt.Printf("  - %s\n", v.Component)
		}

		return fmt.Errorf("sandbox parity violations detected")
	},
}
