package cli

import (
	"context"

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
		repo := InitRepo()

		listings, _, err := repo.FindAll(context.Background(), "", "", "", 0.0, 0.0, 0.0, "", "", false, 100, 0)
		ExitOnErr(err, "Failed to list listings")

		if PrintListResponse(cmd, listings, len(listings), "No listings found") {
			return
		}

		cmd.Printf("Found %d listings:\n\n", len(listings))
		for _, l := range listings {
			PrintListingSummary(cmd, l)
		}
	},
}

var listingGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a listing by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := InitRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		ExitOnErr(err, "Failed to get listing")

		if !FlagText {
			if PrintListResponse(cmd, listing, 1, "") {
				return
			}
		}

		PrintListing(cmd, listing)
	},
}
