package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Generate a context-efficient codebase map",
	RunE: func(cmd *cobra.Command, args []string) error {
		symbols, _ := cmd.Flags().GetBool("symbols")
		templates, _ := cmd.Flags().GetBool("templates")
		depth, _ := cmd.Flags().GetInt("depth")

		if symbols {
			fmt.Println(maintenance.GenerateSymbolMap("."))
			return nil
		}

		if templates {
			fmt.Println(maintenance.GenerateTemplateMap("."))
			return nil
		}

		fmt.Println(maintenance.GeneratePrunedTree(".", depth))
		return nil
	},
}

func init() {
	mapCmd.Flags().Int("depth", 2, "Max depth for directory tree")
	mapCmd.Flags().Bool("symbols", false, "Generate symbol map")
	mapCmd.Flags().Bool("templates", false, "Generate template map")
}
