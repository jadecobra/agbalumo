package cmd

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var listingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all listings",
	Long: `List all existing listings in the agbalumo directory. This command 
supports filtering and can output the results in a machine-readable JSON format.`,
	Example: `  # List all listings
  agbalumo listing list

  # List all listings in JSON format
  agbalumo listing list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listings, _, err := repo.FindAll(context.Background(), "", "", "", "", false, 100, 0)
		if err != nil {
			slog.Error("Failed to list listings", "error", err)
			os.Exit(1)
		}

		if len(listings) == 0 {
			if flagJSON {
				cmd.Println("[]")
			} else {
				cmd.Println("No listings found")
			}
			return
		}

		if flagJSON {
			data, _ := json.MarshalIndent(listings, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Found %d listings:\n\n", len(listings))
		for _, l := range listings {
			printListingSummary(cmd, l)
		}
	},
}

var listingGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a listing by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		if err != nil {
			slog.Error("Failed to get listing", "error", err)
			os.Exit(1)
		}

		if flagJSON {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		printListing(cmd, listing)
	},
}
