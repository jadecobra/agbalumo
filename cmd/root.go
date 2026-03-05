package cmd

import (
	"github.com/spf13/cobra"
)

var flagJSON bool

var rootCmd = &cobra.Command{
	Use:   "agbalumo",
	Short: "agbalumo is a directory and request platform",
	Long: `agbalumo is a high-performance directory and request platform tailored for the 
West African diaspora. It provides a robust command-line interface to manage 
listings, categories, and administrative tasks efficiently.

Core Features:
  - Manage business and service listings.
  - Organize listings into customizable categories.
  - Perform administrative actions like approving/rejecting listings.
  - Seed and serve the platform data.`,
	Example: `  # Start the web server
  agbalumo serve

  # List all listings in JSON format
  agbalumo listing list --json

  # Add a new category
  agbalumo category add "Healthcare" --claimable`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format")
}
