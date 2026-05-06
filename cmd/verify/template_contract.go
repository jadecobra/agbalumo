package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var templateContractCmd = &cobra.Command{
	Use:   "template-contract",
	Short: "Cross-references template field references against ViewModel struct fields",
	RunE: func(cmd *cobra.Command, args []string) error {
		violations, err := maintenance.CheckTemplateContracts(".")
		if err != nil {
			return err
		}

		if len(violations) > 0 {
			fmt.Printf("❌ Found %d template contract violations:\n", len(violations))
			for _, v := range violations {
				fmt.Printf("  - Template: %s, Field: %s, Missing in: %s\n", v.Template, v.Field, v.ExpectedIn)
			}
			os.Exit(1)
		}

		fmt.Println("✅ All templates match ViewModel contracts.")
		return nil
	},
}
