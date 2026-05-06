package cli

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
		repo := InitRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		ExitOnErr(err, "Listing not found")

		ApplyListingUpdates(&listing)

		ExitOnErr(repo.Save(context.Background(), listing), domain.MsgFailedToUpdateListing)

		if !FlagText {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Listing updated successfully: %s\n", listing.ID)
		PrintListing(cmd, listing)
	},
}
