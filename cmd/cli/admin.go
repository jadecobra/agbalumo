package cli

import (
	"context"
	"fmt"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin operations",
	Long: `The admin command provides administrative subcommands for managing 
the agbalumo platform, including approving listings, managing users, 
and viewing claim requests.`,
}

func runListingStatusCmd(status domain.ListingStatus, action string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		repo := InitRepo()
		listing, err := repo.FindByID(context.Background(), args[0])
		ExitOnErr(err, "Listing not found")
		listing.Status = status
		ExitOnErr(repo.Save(context.Background(), listing), fmt.Sprintf("Failed to %s listing", action))
		fmt.Printf("Listing %sd: %s\n", action, args[0])
	}
}

func makeAdminStatusCmd(use, short, action string, status domain.ListingStatus) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		Run:   runListingStatusCmd(status, action),
	}
}

var adminApproveCmd = makeAdminStatusCmd("approve [id]", "Approve a listing", "approve", domain.ListingStatusApproved)

var adminRejectCmd = makeAdminStatusCmd("reject [id]", "Reject a listing", "reject", domain.ListingStatusRejected)

var adminFeaturedCmd = &cobra.Command{
	Use:   "featured [id]",
	Short: "Toggle featured status of a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := InitRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		ExitOnErr(err, "Listing not found")

		listing.Featured = !listing.Featured
		ExitOnErr(repo.Save(context.Background(), listing), domain.MsgFailedToUpdateListing)

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
		repo := InitRepo()

		claims, err := repo.GetPendingClaimRequests(context.Background())
		ExitOnErr(err, "Failed to get pending claims")

		if PrintListResponse(cmd, claims, len(claims), "No pending claim requests") {
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
		repo := InitRepo()

		users, err := repo.GetAllUsers(context.Background(), 100, 0)
		ExitOnErr(err, "Failed to get users")

		if PrintListResponse(cmd, users, len(users), "No users found") {
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
		repo := InitRepo()

		user, err := repo.FindUserByID(context.Background(), args[0])
		ExitOnErr(err, "User not found")

		user.Role = domain.UserRoleAdmin
		ExitOnErr(repo.SaveUser(context.Background(), user), "Failed to promote user")

		fmt.Printf("User promoted to admin: %s\n", args[0])
	},
}

func init() {
	AdminCmd.AddCommand(adminApproveCmd)
	AdminCmd.AddCommand(adminRejectCmd)
	AdminCmd.AddCommand(adminFeaturedCmd)
	AdminCmd.AddCommand(adminPendingClaimsCmd)
	AdminCmd.AddCommand(adminUsersCmd)
	AdminCmd.AddCommand(adminPromoteCmd)

}
