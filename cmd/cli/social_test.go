package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/cmd/cli"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocialDraftCmd(t *testing.T) {
	tempDir := t.TempDir()
	tempDB := filepath.Join(tempDir, "test_social_cli.db")
	_ = os.Setenv("DATABASE_URL", tempDB)
	defer func() {
		_ = os.Unsetenv("DATABASE_URL")
	}()

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

	tests := []struct {
		name     string
		city     string
		platform string
		contains []string
		pillar   int
	}{
		{
			name:     "pillar_1_quality_index_facebook",
			city:     "Dallas",
			platform: "facebook",
			contains: []string{
				"Mama Put Dallas",
				"4.7",
				"128 reviews",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=quality_index&utm_medium=social&utm_source=facebook",
				"https://agbalumo.com/?city=Dallas&utm_campaign=quality_index&utm_medium=social&utm_source=facebook",
				"https://agbalumo.com/?action=post&utm_campaign=quality_index&utm_medium=social&utm_source=facebook",
				"Yoruba",
			},
			pillar: 1,
		},
		{
			name:     "pillar_2_airport_arrival_facebook",
			city:     "Dallas",
			platform: "facebook",
			contains: []string{
				"airport",
				"T's Buka",
				"+16822766400",
				"https://agbalumo.com/listings/test-arlington-1?utm_campaign=airport_arrival&utm_medium=social&utm_source=facebook",
				"https://agbalumo.com/?city=Arlington&utm_campaign=airport_arrival&utm_medium=social&utm_source=facebook",
				"https://agbalumo.com/?action=post&utm_campaign=airport_arrival&utm_medium=social&utm_source=facebook",
			},
			pillar: 2,
		},
		{
			name:     "pillar_3_merchant_spotlight_facebook",
			city:     "Dallas",
			platform: "facebook",
			contains: []string{
				"Spotlight",
				"Mama Put Dallas",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=merchant_spotlight&utm_medium=social&utm_source=facebook",
			},
			pillar: 3,
		},
		{
			name:     "pillar_4_sub_metro_corridor_facebook",
			city:     "Plano",
			platform: "facebook",
			contains: []string{
				"Collin County",
				"BOXOCHOPS",
				"https://agbalumo.com/listings/test-plano-1?utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=facebook",
				"https://agbalumo.com/?city=Plano&utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=facebook",
				"https://agbalumo.com/?action=post&utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=facebook",
			},
			pillar: 4,
		},
		{
			name:     "pillar_5_coverage_gaps_facebook",
			city:     "Dallas",
			platform: "facebook",
			contains: []string{
				"verified",
				"Mama Put Dallas",
				"4.7",
				"128 reviews",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=coverage_gaps&utm_medium=social&utm_source=facebook",
				"Who are we missing",
				"https://agbalumo.com/?action=post&utm_campaign=coverage_gaps&utm_medium=social&utm_source=facebook",
			},
			pillar: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := cli.GenerateSocialDraft(repo, tc.pillar, tc.platform, tc.city, buf)
			require.NoError(t, err)

			output := buf.String()
			for _, exp := range tc.contains {
				assert.Contains(t, output, exp)
			}
		})
	}
}
