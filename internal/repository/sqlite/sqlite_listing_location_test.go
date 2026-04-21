package sqlite_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListingLocationFields(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// This should fail to compile if State/Country are missing
	l := domain.Listing{
		ID:       "loc-test",
		Title:    "Location Test",
		City:     "Dallas",
		State:    "TX",
		Country:  "USA",
		IsActive: true,
		Status:   domain.ListingStatusApproved,
	}

	err := repo.Save(ctx, l)
	assert.NoError(t, err)

	got, err := repo.FindByID(ctx, "loc-test")
	assert.NoError(t, err)
	assert.Equal(t, "TX", got.State)
	assert.Equal(t, "USA", got.Country)
}
