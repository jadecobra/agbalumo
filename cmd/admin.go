package cmd

import (
	"context"
	"encoding/json"
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
	Long: `The admin command provides administrative subcommands for managing 
the agbalumo platform, including approving listings, managing users, 
and viewing claim requests.`,
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

var adminPendingClaimsCmd = &cobra.Command{
	Use:   "pending-claims",
	Short: "List pending claim requests",
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		claims, err := repo.GetPendingClaimRequests(context.Background())
		if err != nil {
			slog.Error("Failed to get pending claims", "error", err)
			os.Exit(1)
		}

		if len(claims) == 0 {
			if flagJSON {
				fmt.Println("[]")
			} else {
				fmt.Println("No pending claim requests")
			}
			return
		}

		if flagJSON {
			data, _ := json.MarshalIndent(claims, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Found %d pending claim requests:\n\n", len(claims))
		for _, cr := range claims {
			cmd.Printf("[%s] Listing: %s | User: %s (%s) | %s\n",
				cr.ID, cr.ListingTitle, cr.UserName, cr.UserEmail, cr.CreatedAt.Format("2006-01-02"))
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
			if flagJSON {
				fmt.Println("[]")
			} else {
				fmt.Println("No users found")
			}
			return
		}

		if flagJSON {
			data, _ := json.MarshalIndent(users, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Found %d users:\n\n", len(users))
		for _, u := range users {
			role := string(u.Role)
			if role == "" {
				role = "user"
			}
			cmd.Printf("[%s] %s - %s\n", u.ID, u.Email, role)
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
	adminCmd.AddCommand(adminPendingClaimsCmd)
	adminCmd.AddCommand(adminUsersCmd)
	adminCmd.AddCommand(adminPromoteCmd)

	rootCmd.AddCommand(adminCmd)
}
