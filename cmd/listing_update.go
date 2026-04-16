package cmd

import (
	"context"
	"encoding/json"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var listingUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		exitOnErr(err, "Listing not found")

		applyListingUpdates(&listing)

		exitOnErr(repo.Save(context.Background(), listing), domain.MsgFailedToUpdateListing)

		if !flagText {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Listing updated successfully: %s\n", listing.ID)
		printListing(cmd, listing)
	},
}
