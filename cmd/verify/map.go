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
		depth, _ := cmd.Flags().GetInt("depth")

		if symbols {
			fmt.Println(maintenance.GenerateSymbolMap("."))
			return nil
		}

		routes, _ := cmd.Flags().GetBool("routes")
		if routes {
			res, err := maintenance.GenerateRouteMap()
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		}

		fmt.Println(maintenance.GeneratePrunedTree(".", depth))
		return nil
	},
}

func init() {
	mapCmd.Flags().Int("depth", 2, "Max depth for directory tree")
	mapCmd.Flags().Bool("symbols", false, "Generate symbol map")
	mapCmd.Flags().Bool("routes", false, "Generate route map from live server bootstrap")
}
