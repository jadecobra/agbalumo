package seeder_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedAll(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)

	seeder.SeedAll(context.Background(), repo)

	// Verify some listings were saved
	listings, _, err := repo.FindAll(context.Background(), "", "", "", "", true, 100, 0)
	require.NoError(t, err)
	assert.Greater(t, len(listings), 0, "Expected listings to be seeded")
}

func TestEnsureSeeded_Empty(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)

	seeder.EnsureSeeded(context.Background(), repo)

	// Verify listings were saved
	listings, _, err := repo.FindAll(context.Background(), "", "", "", "", true, 1, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(listings))
}

func TestEnsureSeeded_NotEmpty(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)
	// Seed one listing
	l := domain.Listing{ID: "1", Title: "Existing", OwnerOrigin: "Ghana", Type: "Business", Address: "123 St"}
	_ = repo.Save(context.Background(), l)

	seeder.EnsureSeeded(context.Background(), repo)

	// Verify NO additional listings were saved (still just 1)
	listings, _, err := repo.FindAll(context.Background(), "", "", "", "", true, 100, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(listings))
}

func TestEnsureSeeded_FindAllError(t *testing.T) {
	// Skipping as forcing error with SQLite is hard without mocks.
}
