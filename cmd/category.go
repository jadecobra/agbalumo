package cmd

import (
	"context"
	"fmt"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage categories",
	Long: `The category command provides subcommands to add and list categories 
used to organize listings in the agbalumo directory.`,
}

var categoryAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new category",
	Long: `Add a new category to the agbalumo system. Categories are used to 
properly classify and filter listings.`,
	Example: `  # Add a new claimable category
  agbalumo category add "Professional Services" --claimable`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		name := args[0]
		claimable, _ := cmd.Flags().GetBool("claimable")

		cat := domain.CategoryData{
			Name:      name,
			Claimable: claimable,
			IsSystem:  false, // user-added are not system categories
			Active:    true,  // active by default
		}

		exitOnErr(repo.SaveCategory(context.Background(), cat), "Failed to save category")

		fmt.Printf("Successfully added category: '%s'\n", name)
	},
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active categories",
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		categories, err := repo.GetCategories(context.Background(), domain.CategoryFilter{ActiveOnly: false})
		exitOnErr(err, "Failed to get categories")

		if printListResponse(cmd, categories, len(categories), "No categories found.") {
			return
		}

		cmd.Printf("\n%-20s %-10s %-15s %-10s\n", "NAME", "ACTIVE", "CLAIMABLE", "SYSTEM")
		cmd.Printf("------------------------------------------------------------\n")
		for _, cat := range categories {
			cmd.Printf("%-20s %-10t %-15t %-10t\n", cat.Name, cat.Active, cat.Claimable, cat.IsSystem)
		}
		cmd.Println()
	},
}

func init() {
	categoryAddCmd.Flags().BoolP("claimable", "c", false, "Is this category claimable?")
	categoryCmd.AddCommand(categoryAddCmd)
	categoryCmd.AddCommand(categoryListCmd)

	rootCmd.AddCommand(categoryCmd)
}
