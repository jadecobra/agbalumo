package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var listingDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := InitRepo()

		if err := repo.Delete(context.Background(), args[0]); err != nil {
			slog.Error("Failed to delete listing", "error", err)
			os.Exit(1)
		}

		fmt.Printf("Listing deleted successfully: %s\n", args[0])
	},
}
