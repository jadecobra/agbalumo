package listing_test

import (
	"context"
	"testing"

	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/require"
)

var testSvcUser = domain.User{ID: "user-123", Name: "Test User", Email: "test@example.com"}

func TestListingService_ClaimListing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success creates pending claim request", func(t *testing.T) {
		t.Parallel()
		repo := testutil.SetupTestRepository(t)
		svc := listmod.NewListingService(repo, repo, repo)

		// Seed listing and category
		saveTestListing(t, repo, "loc-123", "Test Listing")
		_ = repo.SaveCategory(ctx, domain.CategoryData{ID: string(domain.Business), Name: "Business", Claimable: true})

		cr, err := svc.ClaimListing(ctx, testSvcUser, "loc-123")
		require.NoError(t, err)
		require.Equal(t, domain.ClaimStatusPending, cr.Status)
		require.Equal(t, testSvcUser.ID, cr.UserID)
		require.NotEmpty(t, cr.ID)

		// Verify in DB
		saved, err := repo.GetClaimRequestByUserAndListing(ctx, testSvcUser.ID, "loc-123")
		require.NoError(t, err)
		require.Equal(t, cr.ID, saved.ID)
	})

	t.Run("missing user id", func(t *testing.T) {
		t.Parallel()
		repo := testutil.SetupTestRepository(t)
		svc := listmod.NewListingService(repo, repo, repo)

		_, err := svc.ClaimListing(ctx, domain.User{}, "loc-123")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrUserIDRequired)
	})

	t.Run("listing not found", func(t *testing.T) {
		t.Parallel()
		repo := testutil.SetupTestRepository(t)
		svc := listmod.NewListingService(repo, repo, repo)

		_, err := svc.ClaimListing(ctx, testSvcUser, "bad-id")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrListingNotFound)
	})

	t.Run("already owned", func(t *testing.T) {
		t.Parallel()
		repo := testutil.SetupTestRepository(t)
		svc := listmod.NewListingService(repo, repo, repo)

		saveTestListing(t, repo, "loc-123", "Test Listing", func(l *domain.Listing) { l.OwnerID = "someone-else" })

		_, err := svc.ClaimListing(ctx, testSvcUser, "loc-123")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrListingOwned)
	})

	t.Run("unclaimable type", func(t *testing.T) {
		t.Parallel()
		repo := testutil.SetupTestRepository(t)
		svc := listmod.NewListingService(repo, repo, repo)

		saveTestListing(t, repo, "loc-123", "Test Job", func(l *domain.Listing) { l.Type = domain.Job })
		_ = repo.SaveCategory(ctx, domain.CategoryData{ID: string(domain.Job), Name: "Job", Claimable: false})

		_, err := svc.ClaimListing(ctx, testSvcUser, "loc-123")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrListingNotClaimable)
	})

	t.Run("duplicate pending claim rejected", func(t *testing.T) {
		t.Parallel()
		repo := testutil.SetupTestRepository(t)
		svc := listmod.NewListingService(repo, repo, repo)

		saveTestListing(t, repo, "loc-123", "Test Listing")
		_ = repo.SaveCategory(ctx, domain.CategoryData{ID: string(domain.Business), Name: "Business", Claimable: true})
		_ = repo.SaveClaimRequest(ctx, domain.ClaimRequest{ID: "existing", UserID: testSvcUser.ID, ListingID: "loc-123", Status: domain.ClaimStatusPending})

		_, err := svc.ClaimListing(ctx, testSvcUser, "loc-123")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrPendingClaimExists)
	})
}
