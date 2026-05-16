package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var surfaceParityCmd = &cobra.Command{
	Use:   "surface-parity",
	Short: "Verify visual token parity between listing cards and modal details",
	Run: func(cmd *cobra.Command, args []string) {
		violations, err := maintenance.CheckSurfaceParity(".")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if len(violations) > 0 {
			fmt.Println("❌ Surface Parity Violations:")
			for _, v := range violations {
				fmt.Printf("→ %s\n", v)
			}
			os.Exit(1)
		}

		fmt.Println("✅ Surface parity confirmed.")
	},
}
