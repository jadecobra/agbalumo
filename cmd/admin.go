package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var (
	flagStatus string
	flagAction string
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin operations",
}

var adminApproveCmd = &cobra.Command{
	Use:   "approve [id]",
	Short: "Approve a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		if err != nil {
			slog.Error("Listing not found", "error", err)
			os.Exit(1)
		}

		listing.Status = domain.ListingStatusApproved
		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error("Failed to approve listing", "error", err)
			os.Exit(1)
		}

		fmt.Printf("Listing approved: %s\n", args[0])
	},
}

var adminRejectCmd = &cobra.Command{
	Use:   "reject [id]",
	Short: "Reject a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		if err != nil {
			slog.Error("Listing not found", "error", err)
			os.Exit(1)
		}

		listing.Status = domain.ListingStatusRejected
		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error("Failed to reject listing", "error", err)
			os.Exit(1)
		}

		fmt.Printf("Listing rejected: %s\n", args[0])
	},
}

var adminFeaturedCmd = &cobra.Command{
	Use:   "featured [id]",
	Short: "Toggle featured status of a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		if err != nil {
			slog.Error("Listing not found", "error", err)
			os.Exit(1)
		}

		listing.Featured = !listing.Featured
		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error("Failed to update listing", "error", err)
			os.Exit(1)
		}

		status := "featured"
		if !listing.Featured {
			status = "unfeatured"
		}
		fmt.Printf("Listing %s: %s\n", status, args[0])
	},
}

var adminPendingCmd = &cobra.Command{
	Use:   "pending",
	Short: "List pending listings",
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listings, err := repo.GetPendingListings(context.Background(), 100, 0)
		if err != nil {
			slog.Error("Failed to get pending listings", "error", err)
			os.Exit(1)
		}

		if len(listings) == 0 {
			fmt.Println("No pending listings")
			return
		}

		fmt.Printf("Found %d pending listings:\n\n", len(listings))
		for _, l := range listings {
			fmt.Printf("[%s] %s - %s (%s)\n", l.ID, l.Title, l.Type, l.Status)
		}
	},
}

var adminUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		users, err := repo.GetAllUsers(context.Background(), 100, 0)
		if err != nil {
			slog.Error("Failed to get users", "error", err)
			os.Exit(1)
		}

		if len(users) == 0 {
			fmt.Println("No users found")
			return
		}

		fmt.Printf("Found %d users:\n\n", len(users))
		for _, u := range users {
			role := string(u.Role)
			if role == "" {
				role = "user"
			}
			fmt.Printf("[%s] %s - %s\n", u.ID, u.Email, role)
		}
	},
}

var adminPromoteCmd = &cobra.Command{
	Use:   "promote [user-id]",
	Short: "Promote a user to admin by user ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		user, err := repo.FindUserByID(context.Background(), args[0])
		if err != nil {
			slog.Error("User not found", "error", err)
			os.Exit(1)
		}

		user.Role = domain.UserRoleAdmin
		if err := repo.SaveUser(context.Background(), user); err != nil {
			slog.Error("Failed to promote user", "error", err)
			os.Exit(1)
		}

		fmt.Printf("User promoted to admin: %s\n", args[0])
	},
}

func init() {
	adminCmd.AddCommand(adminApproveCmd)
	adminCmd.AddCommand(adminRejectCmd)
	adminCmd.AddCommand(adminFeaturedCmd)
	adminCmd.AddCommand(adminPendingCmd)
	adminCmd.AddCommand(adminUsersCmd)
	adminCmd.AddCommand(adminPromoteCmd)

	rootCmd.AddCommand(adminCmd)
}
