package cli_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/cmd/cli"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) domain.ListingRepository {
	t.Helper()
	tempDir := t.TempDir()
	tempDB := filepath.Join(tempDir, "test_social_cli.db")
	t.Setenv("DATABASE_URL", tempDB)
	t.Setenv("AGBALUMO_SOCIAL_STATE", filepath.Join(tempDir, "social_state.json"))

	repo := cli.InitRepo()
	ctx := context.Background()

	// Seed sample test listings
	_ = repo.Save(ctx, domain.Listing{
		ID:                "test-dallas-1",
		Title:             "Mama Put Dallas",
		Type:              domain.Food,
		City:              "Dallas",
		RegionalSpecialty: "Yoruba",
		Rating:            4.7,
		ReviewCount:       128,
		ContactPhone:      "+12145551001",
		WebsiteURL:        "https://mamaputdallas.com",
		Status:            domain.ListingStatusApproved,
		IsActive:          true,
	})
	_ = repo.Save(ctx, domain.Listing{
		ID:                "test-plano-1",
		Title:             "BOXOCHOPS",
		Type:              domain.Food,
		City:              "Plano",
		RegionalSpecialty: "Nigerian",
		Rating:            4.6,
		ReviewCount:       178,
		ContactPhone:      "+14692940003",
		WebsiteURL:        "https://toasttab.com/boxochops",
		Status:            domain.ListingStatusApproved,
		IsActive:          true,
	})
	_ = repo.Save(ctx, domain.Listing{
		ID:                "test-arlington-1",
		Title:             "T's Buka",
		Type:              domain.Food,
		City:              "Arlington",
		RegionalSpecialty: "Nigerian",
		Rating:            4.1,
		ReviewCount:       89,
		ContactPhone:      "+16822766400",
		Status:            domain.ListingStatusApproved,
		IsActive:          true,
	})

	return repo
}

func TestSocialDraftCmd(t *testing.T) {
	repo := setupTestRepo(t)

	tests := []struct {
		name      string
		city      string
		listingID string
		contains  []string
		omits     []string
		pillar    int
	}{
		{
			name:      "pillar_1_quality_index_honest_utm",
			city:      "Dallas",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 1: The DFW African Food Quality Index]",
				"Mama Put Dallas",
				"4.7",
				"128 reviews",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=quality_index&utm_medium=social&utm_source=cli",
				"https://agbalumo.com/?city=Dallas&utm_campaign=quality_index&utm_medium=social&utm_source=cli",
				"https://agbalumo.com/?action=post&utm_campaign=quality_index&utm_medium=social&utm_source=cli",
				"Yoruba",
			},
			omits:  nil,
			pillar: 1,
		},
		{
			name:      "pillar_2_airport_corridor_honest_copy",
			city:      "Dallas",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 2: Airport Corridor Cities (Arlington, Grand Prairie, Irving)]",
				"T's Buka",
				"+16822766400",
				"https://agbalumo.com/listings/test-arlington-1?utm_campaign=airport_corridor&utm_medium=social&utm_source=cli",
				"https://agbalumo.com/?city=Arlington&utm_campaign=airport_corridor&utm_medium=social&utm_source=cli",
				"https://agbalumo.com/?action=post&utm_campaign=airport_corridor&utm_medium=social&utm_source=cli",
			},
			omits: []string{
				"late-night",
				"kitchen is still open",
				"within 20 minutes",
				"after 8 PM",
			},
			pillar: 2,
		},
		{
			name:      "pillar_3_merchant_spotlight_specified_id",
			city:      "Dallas",
			listingID: "test-dallas-1",
			contains: []string{
				"[DRAFT - Pillar 3: Merchant Spotlight]",
				"Mama Put Dallas",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=merchant_spotlight&utm_medium=social&utm_source=cli",
			},
			omits:  nil,
			pillar: 3,
		},
		{
			name:      "pillar_4_sub_metro_corridor_honest_copy",
			city:      "Plano",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 4: Sub-Metro Corridor Guide (Collin County)]",
				"Collin County",
				"BOXOCHOPS",
				"https://agbalumo.com/listings/test-plano-1?utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=cli",
				"https://agbalumo.com/?city=Plano&utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=cli",
				"https://agbalumo.com/?action=post&utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=cli",
			},
			omits: []string{
				"45 minutes",
			},
			pillar: 4,
		},
		{
			name:      "pillar_5_coverage_gaps",
			city:      "Dallas",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 5: Radical Transparency & Coverage Gaps]",
				"verified",
				"Mama Put Dallas",
				"4.7",
				"128 reviews",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=coverage_gaps&utm_medium=social&utm_source=cli",
				"Who are we missing",
				"https://agbalumo.com/?action=post&utm_campaign=coverage_gaps&utm_medium=social&utm_source=cli",
			},
			omits:  nil,
			pillar: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := cli.GenerateSocialDraft(repo, tc.pillar, tc.listingID, tc.city, buf)
			require.NoError(t, err)

			output := buf.String()
			for _, exp := range tc.contains {
				assert.Contains(t, output, exp)
			}
			for _, bad := range tc.omits {
				assert.NotContains(t, output, bad)
			}
		})
	}
}

func TestSocialDraft_Pillar3NotFound(t *testing.T) {
	repo := setupTestRepo(t)
	buf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 3, "non-existent-id", "Dallas", buf)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-existent-id")
}

func TestSocialDraft_Pillar3RotationWithoutID(t *testing.T) {
	repo := setupTestRepo(t)

	// Call pillar 3 twice without listing-id. Should spotlight different listings across runs.
	buf1 := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 3, "", "Dallas", buf1)
	require.NoError(t, err)

	buf2 := new(bytes.Buffer)
	err = cli.GenerateSocialDraft(repo, 3, "", "Dallas", buf2)
	require.NoError(t, err)

	require.NotEqual(t, buf1.String(), buf2.String(), "Consecutive Pillar 3 drafts without listing-id should rotate spotlighted merchant")
}

func TestSocialDraft_RotationAcrossListings(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	// Add 5 more Dallas listings to have >4 listings for rotation testing
	extraListings := []domain.Listing{
		{ID: "dallas-extra-1", Title: "Extra Spot 1", Type: domain.Food, City: "Dallas", Rating: 4.8, ReviewCount: 50, Status: domain.ListingStatusApproved, IsActive: true},
		{ID: "dallas-extra-2", Title: "Extra Spot 2", Type: domain.Food, City: "Dallas", Rating: 4.5, ReviewCount: 40, Status: domain.ListingStatusApproved, IsActive: true},
		{ID: "dallas-extra-3", Title: "Extra Spot 3", Type: domain.Food, City: "Dallas", Rating: 4.4, ReviewCount: 30, Status: domain.ListingStatusApproved, IsActive: true},
		{ID: "dallas-extra-4", Title: "Extra Spot 4", Type: domain.Food, City: "Dallas", Rating: 4.2, ReviewCount: 20, Status: domain.ListingStatusApproved, IsActive: true},
		{ID: "dallas-extra-5", Title: "Extra Spot 5", Type: domain.Food, City: "Dallas", Rating: 4.0, ReviewCount: 10, Status: domain.ListingStatusApproved, IsActive: true},
	}
	for _, l := range extraListings {
		_ = repo.Save(ctx, l)
	}

	buf1 := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 1, "", "Dallas", buf1)
	require.NoError(t, err)

	buf2 := new(bytes.Buffer)
	err = cli.GenerateSocialDraft(repo, 1, "", "Dallas", buf2)
	require.NoError(t, err)

	assert.NotEqual(t, buf1.String(), buf2.String(), "Pillar 1 should rotate listings across calls so different spots get exposure")
}
