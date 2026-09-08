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
			pillar:   1,
			platform: "facebook",
			city:     "Dallas",
			contains: []string{
				"Mama Put Dallas",
				"4.7",
				"https://agbalumo.com/?city=Dallas",
				"Yoruba",
			},
		},
		{
			name:     "pillar_2_airport_arrival_facebook",
			pillar:   2,
			platform: "facebook",
			city:     "Dallas",
			contains: []string{
				"airport",
				"T's Buka",
				"+16822766400",
				"https://agbalumo.com/?city=Arlington",
			},
		},
		{
			name:     "pillar_3_merchant_spotlight_facebook",
			pillar:   3,
			platform: "facebook",
			city:     "Dallas",
			contains: []string{
				"Spotlight",
				"Mama Put Dallas",
				"https://agbalumo.com/listings/test-dallas-1",
			},
		},
		{
			name:     "pillar_4_sub_metro_corridor_facebook",
			pillar:   4,
			platform: "facebook",
			city:     "Plano",
			contains: []string{
				"Collin County",
				"BOXOCHOPS",
				"https://agbalumo.com/?city=Plano",
			},
		},
		{
			name:     "pillar_5_coverage_gaps_facebook",
			pillar:   5,
			platform: "facebook",
			city:     "Dallas",
			contains: []string{
				"verified",
				"Mama Put Dallas",
				"Who are we missing",
				"https://agbalumo.com",
			},
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
