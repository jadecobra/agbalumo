package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/cmd/cli"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepoWithState(t *testing.T) (domain.ListingRepository, string) {
	t.Helper()
	tempDir := t.TempDir()
	tempDB := filepath.Join(tempDir, "test_social_cli.db")
	t.Setenv("DATABASE_URL", tempDB)
	t.Setenv("AGBALUMO_SOCIAL_DB", tempDB)
	t.Setenv("AGBALUMO_SKIP_LINK_VERIFY", "true")
	statePath := filepath.Join(tempDir, "social_state.json")
	t.Setenv("AGBALUMO_SOCIAL_STATE", statePath)

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

	return repo, statePath
}

func setupTestRepo(t *testing.T) domain.ListingRepository {
	t.Helper()
	repo, _ := setupTestRepoWithState(t)
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
				"Landing in a new city or moving across town shouldn't mean gambling on food quality.",
				"Here are places in DFW where quality is backed by real Google reviews:",
				"Mama Put Dallas",
				"4.7",
				"128 reviews",
				"Reviews & Details: https://agbalumo.com/listings/test-dallas-1?utm_campaign=quality_index&utm_medium=social&utm_source=facebook",
				"Phone: +12145551001",
				"Menu/Order: https://mamaputdallas.com",
				"Find what you want in under 60 seconds:",
				"https://agbalumo.com/?utm_campaign=quality_index&utm_medium=social&utm_source=facebook",
				"If we missed places you like in DFW, add it directly to the network:",
				"https://agbalumo.com/?action=post&utm_campaign=quality_index&utm_medium=social&utm_source=facebook",
				"Yoruba",
			},
			omits: []string{
				"community reviews",
				"add it directly to the network in under 60 seconds",
			},
			pillar: 1,
		},
		{
			name:      "pillar_2_airport_corridor_honest_copy",
			city:      "Dallas",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 2: Airport Corridor Cities (Arlington, Grand Prairie, Irving)]",
				"You don't have to drive or search the entire metroplex to find African food when you land at DFW.",
				"Here are verified food places in the Arlington, Grand Prairie, Irving area with direct contact information: The reviews are from google",
				"T's Buka",
				"Reviews & Details: https://agbalumo.com/listings/test-arlington-1?utm_campaign=airport_corridor&utm_medium=social&utm_source=facebook",
				"Direct Phone: +16822766400",
				"Explore more African food in the DFW airport corridor at",
				"https://agbalumo.com/?utm_campaign=airport_corridor&utm_medium=social&utm_source=facebook",
				"If you know other African-owned places near DFW airport that we missed, add it directly at:",
				"https://agbalumo.com/?action=post&utm_campaign=airport_corridor&utm_medium=social&utm_source=facebook",
			},
			omits: []string{
				"late-night",
				"kitchen is still open",
				"within 20 minutes",
				"after 8 PM",
				"community reviews",
				"under 60 seconds",
			},
			pillar: 2,
		},
		{
			name:      "pillar_3_merchant_spotlight_specified_id",
			city:      "Dallas",
			listingID: "test-dallas-1",
			contains: []string{
				"[DRAFT - Pillar 3: Merchant Spotlight]",
				"Have you been to Mama Put Dallas in Dallas",
				"they have a ★ 4.7 rating across 128 Google reviews.",
				"+12145551001",
				"https://mamaputdallas.com",
				"See more at https://agbalumo.com/listings/test-dallas-1?utm_campaign=merchant_spotlight&utm_medium=social&utm_source=facebook",
			},
			omits: []string{
				"Spotlight:",
				"Serving authentic",
				"Tagging",
				"thank you for serving",
				"community reviews",
			},
			pillar: 3,
		},
		{
			name:      "pillar_4_sub_metro_corridor_honest_copy",
			city:      "Plano",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 4: Sub-Metro Corridor Guide (Collin County)]",
				"Craving authentic West African food?",
				"There is a trusted cluster of verified African kitchens with Google Reviews in Plano, Allen, McKinney, and Frisco:",
				"BOXOCHOPS",
				"Reviews & Menu: https://agbalumo.com/listings/test-plano-1?utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=facebook",
				"Phone: +14692940003",
				"Online Order: https://toasttab.com/boxochops",
				"Explore North DFW and Collin County food places at",
				"https://agbalumo.com/?utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=facebook",
				"Know another African-owned kitchen in Collin County? Add it directly to the network:",
				"https://agbalumo.com/?action=post&utm_campaign=sub_metro_corridor&utm_medium=social&utm_source=facebook",
			},
			omits: []string{
				"45 minutes",
				"community reviews",
				"under 60 seconds",
				"Direct Phone",
			},
			pillar: 4,
		},
		{
			name:      "pillar_5_coverage_gaps",
			city:      "Dallas",
			listingID: "",
			contains: []string{
				"[DRAFT - Pillar 5: Radical Transparency & Coverage Gaps]",
				"Here are places with Google reviews:",
				"Mama Put Dallas",
				"4.7",
				"128 reviews",
				"https://agbalumo.com/listings/test-dallas-1?utm_campaign=coverage_gaps&utm_medium=social&utm_source=facebook",
				"Find other places at",
				"https://agbalumo.com/?utm_campaign=coverage_gaps&utm_medium=social&utm_source=facebook",
				"Who are we missing? Add a place you like:",
				"https://agbalumo.com/?action=post&utm_campaign=coverage_gaps&utm_medium=social&utm_source=facebook",
			},
			omits: []string{
				"community reviews",
				"under 60 seconds",
				"favorite auntie's spot",
				"suya joint",
			},
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

func TestSocialDraft_Pillar3UnknownListingLeavesStateUnchanged(t *testing.T) {
	repo, statePath := setupTestRepoWithState(t)

	// 1. Initial state: state file must not exist yet.
	_, err := os.Stat(statePath)
	require.True(t, os.IsNotExist(err))

	buf := new(bytes.Buffer)
	err = cli.GenerateSocialDraft(repo, 3, "unknown-listing-id", "Dallas", buf)
	require.Error(t, err)

	// When generation fails, no state file should be created.
	_, err = os.Stat(statePath)
	assert.True(t, os.IsNotExist(err), "state file should not be created when draft generation fails")

	// 2. Pre-existing state: state file has existing offsets and recently_featured.
	initialJSON := "{\n  \"offsets\": {\n    \"pillar_1_dallas\": 3,\n    \"pillar_3\": 1\n  },\n  \"recently_featured\": [\n    \"test-dallas-1\"\n  ]\n}"
	require.NoError(t, os.WriteFile(statePath, []byte(initialJSON), 0600))
	infoBefore, err := os.Stat(statePath)
	require.NoError(t, err)

	buf.Reset()
	err = cli.GenerateSocialDraft(repo, 3, "unknown-listing-id", "Dallas", buf)
	require.Error(t, err)

	infoAfter, err := os.Stat(statePath)
	require.NoError(t, err)
	assert.Equal(t, infoBefore.ModTime(), infoAfter.ModTime(), "state file should not be touched on failed draft")

	contentAfter, err := os.ReadFile(filepath.Clean(statePath)) // #nosec G304 -- test reading local state path
	require.NoError(t, err)
	assert.Equal(t, initialJSON, string(contentAfter), "offsets and recently_featured should be unchanged")
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

func TestSocialDraft_RotateDistinctCities(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	// Seed multiple listings across cities with some repeats:
	// Dallas (3 spots), Fort Worth (1 spot), Plano (1 spot), Arlington (already has 1)
	listings := []domain.Listing{
		{ID: "dallas-repeat-1", Title: "Dallas Spot A", Type: domain.Food, City: "Dallas", Rating: 4.8, ReviewCount: 50, Status: domain.ListingStatusApproved, IsActive: true},
		{ID: "dallas-repeat-2", Title: "Dallas Spot B", Type: domain.Food, City: "Dallas", Rating: 4.7, ReviewCount: 40, Status: domain.ListingStatusApproved, IsActive: true},
		{ID: "fw-distinct-1", Title: "FW Spot A", Type: domain.Food, City: "Fort Worth", Rating: 4.6, ReviewCount: 30, Status: domain.ListingStatusApproved, IsActive: true},
	}
	for _, l := range listings {
		_ = repo.Save(ctx, l)
	}

	buf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 1, "", "Dallas", buf)
	require.NoError(t, err)

	out := buf.String()
	// Total spots selected for Pillar 1 is 4.
	// Since we have Dallas, Plano, Arlington, and Fort Worth, all 4 cities must be represented (max 1 per city when possible).
	assert.Contains(t, out, "(Dallas)")
	assert.Contains(t, out, "(Plano)")
	assert.Contains(t, out, "(Arlington)")
	assert.Contains(t, out, "(Fort Worth)")
}

func TestGetSocialDatabaseURL(t *testing.T) {
	tempDir := t.TempDir()
	validDB := filepath.Join(tempDir, "prod_snapshot.db")
	require.NoError(t, os.WriteFile(validDB, []byte("sqlite"), 0600))

	t.Run("rejects_tester_agbalumo_db_via_social_env", func(t *testing.T) {
		t.Setenv("AGBALUMO_SOCIAL_DB", ".tester/data/agbalumo.db")
		_, err := cli.GetSocialDatabaseURL()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "prohibited")
	})

	t.Run("rejects_tester_agbalumo_db_via_database_url", func(t *testing.T) {
		t.Setenv("AGBALUMO_SOCIAL_DB", "")
		t.Setenv("DATABASE_URL", ".tester/data/agbalumo.db")
		_, err := cli.GetSocialDatabaseURL()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "prohibited")
	})

	t.Run("prefers_agbalumo_social_db_over_database_url", func(t *testing.T) {
		fallbackDB := filepath.Join(tempDir, "fallback.db")
		t.Setenv("AGBALUMO_SOCIAL_DB", validDB)
		t.Setenv("DATABASE_URL", fallbackDB)

		res, err := cli.GetSocialDatabaseURL()
		require.NoError(t, err)
		assert.Equal(t, validDB, res)
	})

	t.Run("falls_back_to_database_url_when_social_db_unset", func(t *testing.T) {
		t.Setenv("AGBALUMO_SOCIAL_DB", "")
		t.Setenv("DATABASE_URL", validDB)

		res, err := cli.GetSocialDatabaseURL()
		require.NoError(t, err)
		assert.Equal(t, validDB, res)
	})
}

func TestSocialDraft_DeepLinkVerification(t *testing.T) {
	repo := setupTestRepo(t)

	t.Run("known_good_id_200_emits_link_pillar_3", func(t *testing.T) {
		buf := new(bytes.Buffer)
		mockVerifier := cli.LinkVerifierFunc(func(ctx context.Context, id string, _ ...string) (bool, int, error) {
			if id == "test-dallas-1" {
				return true, 200, nil
			}
			return false, 404, nil
		})

		err := cli.GenerateSocialDraft(repo, 3, "test-dallas-1", "Dallas", buf, cli.WithLinkVerifier(mockVerifier))
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "https://agbalumo.com/listings/test-dallas-1?utm_campaign=merchant_spotlight")
	})

	t.Run("wrong_id_404_fails_draft_pillar_3", func(t *testing.T) {
		buf := new(bytes.Buffer)
		mockVerifier := cli.LinkVerifierFunc(func(ctx context.Context, id string, _ ...string) (bool, int, error) {
			return false, 404, nil
		})

		err := cli.GenerateSocialDraft(repo, 3, "test-dallas-1", "Dallas", buf, cli.WithLinkVerifier(mockVerifier))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "deep link verification failed")
	})

	t.Run("wrong_id_404_omits_link_with_loud_stderr_pillar_1", func(t *testing.T) {
		buf := new(bytes.Buffer)
		stderrBuf := new(bytes.Buffer)

		mockVerifier := cli.LinkVerifierFunc(func(ctx context.Context, id string, _ ...string) (bool, int, error) {
			if id == "test-dallas-1" {
				return false, 404, nil
			}
			return true, 200, nil
		})

		err := cli.GenerateSocialDraft(repo, 1, "", "Dallas", buf,
			cli.WithLinkVerifier(mockVerifier),
			cli.WithStderr(stderrBuf),
		)
		require.NoError(t, err)

		// Mama Put Dallas (test-dallas-1) must still appear as a spot
		assert.Contains(t, buf.String(), "Mama Put Dallas")
		// But its deep link must be omitted!
		assert.NotContains(t, buf.String(), "https://agbalumo.com/listings/test-dallas-1")
		// Loud stderr must be recorded!
		assert.Contains(t, stderrBuf.String(), "WARNING")
		assert.Contains(t, stderrBuf.String(), "test-dallas-1")
	})

	t.Run("fail_bad_links_flag_fails_pillar_1", func(t *testing.T) {
		buf := new(bytes.Buffer)
		mockVerifier := cli.LinkVerifierFunc(func(ctx context.Context, id string, _ ...string) (bool, int, error) {
			return false, 404, nil
		})

		err := cli.GenerateSocialDraft(repo, 1, "", "Dallas", buf,
			cli.WithLinkVerifier(mockVerifier),
			cli.WithFailBadLinks(true),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "deep link verification failed")
	})
}

func TestHTTPLinkVerifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		expectedTitle string
		statusCode    int
		wantStatus    int
		wantOk        bool
		wantErr       bool
	}{
		{
			name:          "200_with_html_and_matching_title_passes",
			statusCode:    http.StatusOK,
			body:          "<!DOCTYPE html><html><head><title>Lagos Kitchen | agbalumo</title></head><body><h1>Lagos Kitchen</h1></body></html>",
			expectedTitle: "Lagos Kitchen",
			wantOk:        true,
			wantStatus:    http.StatusOK,
			wantErr:       false,
		},
		{
			name:          "200_missing_html_fails",
			statusCode:    http.StatusOK,
			body:          `{"status":"ok","title":"Lagos Kitchen"}`,
			expectedTitle: "Lagos Kitchen",
			wantOk:        false,
			wantStatus:    http.StatusOK,
			wantErr:       true,
		},
		{
			name:          "200_with_html_but_missing_title_fails",
			statusCode:    http.StatusOK,
			body:          "<!DOCTYPE html><html><head><title>Something Else</title></head><body></body></html>",
			expectedTitle: "Lagos Kitchen",
			wantOk:        false,
			wantStatus:    http.StatusOK,
			wantErr:       true,
		},
		{
			name:          "404_status_fails",
			statusCode:    http.StatusNotFound,
			body:          "<!DOCTYPE html><html><body>Not Found</body></html>",
			expectedTitle: "Lagos Kitchen",
			wantOk:        false,
			wantStatus:    http.StatusNotFound,
			wantErr:       false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			verifier := cli.NewHTTPLinkVerifier(server.URL, server.Client())
			ok, status, err := verifier.Verify(context.Background(), "test-id", tc.expectedTitle)
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.wantStatus, status)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHTTPLinkVerifier_HTMLEntitiesInTitlePasses(t *testing.T) {
	t.Skip("RED: html entities in title")
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><head><title>Mabel&#39;s African Cuisine &amp; Bar | agbalumo</title></head><body><h1>Mabel&#39;s African Cuisine &amp; Bar</h1></body></html>"))
	}))
	defer server.Close()

	verifier := cli.NewHTTPLinkVerifier(server.URL, server.Client())
	ok, status, err := verifier.Verify(context.Background(), "test-id", "Mabel's African Cuisine & Bar")
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, status)
	assert.NoError(t, err)
}

func TestHTTPLinkVerifier_HangsOnFirstRequestRetriesAndPasses(t *testing.T) {
	t.Skip("RED: retry on timeout")
	t.Parallel()

	var reqCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&reqCount, 1)
		if count == 1 {
			select {
			case <-r.Context().Done():
			case <-time.After(500 * time.Millisecond):
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><head><title>Mabel's African Cuisine & Bar</title></head><body></body></html>"))
	}))
	defer server.Close()

	client := server.Client()
	client.Timeout = 50 * time.Millisecond

	verifier := cli.NewHTTPLinkVerifier(server.URL, client)

	ok, status, err := verifier.Verify(context.Background(), "test-id", "Mabel's African Cuisine & Bar")
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, status)
	assert.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
}

func TestHTTPLinkVerifier_NotFoundDoesNotRetry(t *testing.T) {
	t.Parallel()

	var reqCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body>Not Found</body></html>"))
	}))
	defer server.Close()

	verifier := cli.NewHTTPLinkVerifier(server.URL, server.Client())
	ok, status, err := verifier.Verify(context.Background(), "test-id", "Mabel's African Cuisine & Bar")
	assert.False(t, ok)
	assert.Equal(t, http.StatusNotFound, status)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount))
}

func TestHTTPLinkVerifier_MissingTitleDoesNotRetry(t *testing.T) {
	t.Parallel()

	var reqCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body><h1>Other Place</h1></body></html>"))
	}))
	defer server.Close()

	verifier := cli.NewHTTPLinkVerifier(server.URL, server.Client())
	ok, status, err := verifier.Verify(context.Background(), "test-id", "Mabel's African Cuisine & Bar")
	assert.False(t, ok)
	assert.Equal(t, http.StatusOK, status)
	assert.Error(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount))
}

func TestSocialDraft_RegionalExploreLinksContainNoCity(t *testing.T) {
	repo := setupTestRepo(t)

	tests := []struct {
		name     string
		city     string
		campaign string
		pillar   int
	}{
		{
			name:     "pillar_1_quality_index",
			city:     "Dallas",
			campaign: "quality_index",
			pillar:   1,
		},
		{
			name:     "pillar_2_airport_corridor",
			city:     "Dallas",
			campaign: "airport_corridor",
			pillar:   2,
		},
		{
			name:     "pillar_4_sub_metro_corridor",
			city:     "Plano",
			campaign: "sub_metro_corridor",
			pillar:   4,
		},
		{
			name:     "pillar_5_coverage_gaps",
			city:     "Dallas",
			campaign: "coverage_gaps",
			pillar:   5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := cli.GenerateSocialDraft(repo, tc.pillar, "", tc.city, buf)
			require.NoError(t, err)

			output := buf.String()
			expectedExploreURL := "https://agbalumo.com/?utm_campaign=" + tc.campaign + "&utm_medium=social&utm_source=facebook"
			assert.Contains(t, output, expectedExploreURL)
			assert.NotContains(t, output, "city=")
		})
	}
}

func TestSocialDraft_Pillar5CoverageGaps_OneSpotPerCityAndCount(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	// Seed Arlington with 2 additional listings so it has 3 listings total
	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-arlington-2",
		Title:       "Arlington Suya Spot",
		Type:        domain.Food,
		City:        "Arlington",
		Rating:      4.9,
		ReviewCount: 150,
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})
	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-arlington-3",
		Title:       "Arlington Jollof Joint",
		Type:        domain.Food,
		City:        "Arlington",
		Rating:      4.5,
		ReviewCount: 50,
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})

	buf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 5, "", "Dallas", buf)
	require.NoError(t, err)

	out := buf.String()
	// Arlington has 3 listings: should show top-rated spot (Arlington Suya Spot) and "+2 more in Arlington"
	assert.Contains(t, out, "Arlington Suya Spot")
	assert.Contains(t, out, "+2 more in Arlington")
	// Arlington's other 2 spots should NOT have deep links in the output
	assert.NotContains(t, out, "https://agbalumo.com/listings/test-arlington-1")
	assert.NotContains(t, out, "https://agbalumo.com/listings/test-arlington-3")
}

func TestSocialDraft_Pillar5CoverageGaps_NoBlindSpotWithTwoOrMoreListings(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	// Seed Fort Worth with 2 listings
	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-fw-1",
		Title:       "Fort Worth Spot 1",
		Type:        domain.Food,
		City:        "Fort Worth",
		Rating:      4.8,
		ReviewCount: 80,
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})
	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-fw-2",
		Title:       "Fort Worth Spot 2",
		Type:        domain.Food,
		City:        "Fort Worth",
		Rating:      4.2,
		ReviewCount: 30,
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})

	buf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 5, "", "Dallas", buf)
	require.NoError(t, err)

	out := buf.String()
	// Fort Worth has 2 listings, so it is an active city and MUST NOT be named in missing places
	assert.Contains(t, out, "• Fort Worth:")
	assert.NotContains(t, out, "missing places in Fort Worth")
	assert.NotContains(t, out, "Fort Worth.")
	// But watchlist cities with 0 listings (like Frisco, Garland, Denton) should still be in missing places
	assert.Contains(t, out, "Frisco")
}

func TestSocialDraft_Pillar5CoverageGaps_CapFoldsExtraCities(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	// Existing repo has Dallas (1), Plano (1), Arlington (1).
	// Add 7 more cities to make 10 cities total:
	// Grand Prairie, Irving, Fort Worth, McKinney, Allen, Frisco, Garland
	moreCities := []string{"Grand Prairie", "Irving", "Fort Worth", "McKinney", "Allen", "Frisco", "Garland"}
	for i, c := range moreCities {
		_ = repo.Save(ctx, domain.Listing{
			ID:          fmt.Sprintf("test-city-%d", i),
			Title:       fmt.Sprintf("%s Spot", c),
			Type:        domain.Food,
			City:        c,
			Rating:      4.5,
			ReviewCount: 20,
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
		})
	}

	buf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 5, "", "", buf)
	require.NoError(t, err)

	out := buf.String()
	// At most 8 cities shown as full sections
	assert.Contains(t, out, "Also mapped in")
	// Total "• " city headers should be at most 8
	countHeaders := strings.Count(out, "• ")
	assert.LessOrEqual(t, countHeaders, 8)
}

func TestSocialDraft_Pillar5CoverageGaps_AllWatchlistCoveredDropsBlindSpotsLine(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	// Fixed watchlist: Frisco, Garland, Denton, Mesquite, Richardson, Carrollton, Lewisville, Fort Worth
	watchlist := []string{
		"Frisco", "Garland", "Denton", "Mesquite", "Richardson", "Carrollton", "Lewisville", "Fort Worth",
	}
	// Seed each with 2 listings so all have >= 2 listings
	for _, c := range watchlist {
		for j := 1; j <= 2; j++ {
			_ = repo.Save(ctx, domain.Listing{
				ID:          fmt.Sprintf("test-%s-%d", strings.ToLower(c), j),
				Title:       fmt.Sprintf("%s Kitchen %d", c, j),
				Type:        domain.Food,
				City:        c,
				Rating:      4.6,
				ReviewCount: 50,
				Status:      domain.ListingStatusApproved,
				IsActive:    true,
			})
		}
	}

	buf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 5, "", "", buf)
	require.NoError(t, err)

	out := buf.String()
	// When all watchlist cities have >= 2 listings, the missing places line must be dropped completely
	assert.NotContains(t, out, "missing places in")
}

func TestSocialDraft_Lint_MessyMenuLinks(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-shopify-1",
		Title:       "LeAnna Chop Grill",
		Type:        domain.Food,
		City:        "Fort Worth",
		Rating:      4.9,
		ReviewCount: 200,
		WebsiteURL:  "https://g2igzc-as.myshopify.com/",
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})

	stderrBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 3, "test-shopify-1", "Fort Worth", outBuf, cli.WithStderr(stderrBuf))
	require.NoError(t, err)

	stderr := stderrBuf.String()
	assert.Contains(t, stderr, "lint: menu link is raw myshopify — LeAnna Chop Grill (test-shopify-1)")
	assert.Contains(t, outBuf.String(), "LeAnna Chop Grill")
}

func TestSocialDraft_Lint_TrackingParamsCleaned(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-tracking-1",
		Title:       "Kanny's Restaurant",
		Type:        domain.Food,
		City:        "Dallas",
		Rating:      4.9,
		ReviewCount: 250,
		WebsiteURL:  "https://kannysrestaurant.com/catering-services/?v=28886f13f578",
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})

	stderrBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 3, "test-tracking-1", "Dallas", outBuf, cli.WithStderr(stderrBuf))
	require.NoError(t, err)

	stderr := stderrBuf.String()
	assert.Contains(t, stderr, "lint: menu link carries tracking params — Kanny's Restaurant (test-tracking-1) (cleaned: https://kannysrestaurant.com/catering-services/)")
}

func TestSocialDraft_Lint_CityListedAndBlindSpot(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-fw-single",
		Title:       "Fort Worth Suya Spot",
		Type:        domain.Food,
		City:        "Fort Worth",
		Rating:      4.9,
		ReviewCount: 150,
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})

	stderrBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 5, "", "", outBuf, cli.WithStderr(stderrBuf))
	require.NoError(t, err)

	stderr := stderrBuf.String()
	assert.Contains(t, stderr, "lint: city Fort Worth is both listed and called a blind spot")
}

func TestSocialDraft_Lint_RepeatFeatures(t *testing.T) {
	repo, statePath := setupTestRepoWithState(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-repeat-1",
		Title:       "Recent Spot",
		Type:        domain.Food,
		City:        "Dallas",
		Rating:      4.9,
		ReviewCount: 100,
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
	})

	yesterday := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	initialJSON := fmt.Sprintf(`{
  "offsets": {},
  "recently_featured": [
    {"id": "test-repeat-1", "featured_at": "%s"}
  ]
}`, yesterday)
	require.NoError(t, os.WriteFile(statePath, []byte(initialJSON), 0600))

	stderrBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 3, "test-repeat-1", "Dallas", outBuf, cli.WithStderr(stderrBuf))
	require.NoError(t, err)

	stderr := stderrBuf.String()
	assert.Contains(t, stderr, "lint: spot test-repeat-1 was featured in the last 3 days")
}

func TestSocialDraft_Lint_TooLong(t *testing.T) {
	repo, _ := setupTestRepoWithState(t)
	ctx := context.Background()

	for i := 1; i <= 12; i++ {
		_ = repo.Save(ctx, domain.Listing{
			ID:          fmt.Sprintf("test-long-%02d", i),
			Title:       fmt.Sprintf("Long Spot %02d", i),
			Type:        domain.Food,
			City:        fmt.Sprintf("City%02d", i),
			Rating:      4.9,
			ReviewCount: 100,
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
		})
	}

	stderrBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	err := cli.GenerateSocialDraft(repo, 5, "", "", outBuf, cli.WithStderr(stderrBuf))
	require.NoError(t, err)

	stderr := stderrBuf.String()
	assert.True(t, strings.Contains(stderr, "lint: draft has more than 10 listing links") || strings.Contains(stderr, "lint: draft body exceeds 2000 characters"))
}

func TestSocialDraft_Lint_StrictModeFailsAndPreservesState(t *testing.T) {
	repo, statePath := setupTestRepoWithState(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{
		ID:         "test-shopify-strict",
		Title:      "LeAnna Chop Grill",
		Type:       domain.Food,
		City:       "Fort Worth",
		Rating:     4.9,
		WebsiteURL: "https://g2igzc-as.myshopify.com/",
		Status:     domain.ListingStatusApproved,
		IsActive:   true,
	})

	initialJSON := "{\n  \"offsets\": {\n    \"pillar_3\": 0\n  }\n}"
	require.NoError(t, os.WriteFile(statePath, []byte(initialJSON), 0600))
	infoBefore, err := os.Stat(statePath)
	require.NoError(t, err)

	stderrBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	err = cli.GenerateSocialDraft(repo, 3, "test-shopify-strict", "Fort Worth", outBuf, cli.WithStderr(stderrBuf), cli.WithStrictMode(true))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lint check failed")

	infoAfter, err := os.Stat(statePath)
	require.NoError(t, err)
	assert.Equal(t, infoBefore.ModTime(), infoAfter.ModTime())

	contentAfter, err := os.ReadFile(filepath.Clean(statePath))
	require.NoError(t, err)
	assert.Equal(t, initialJSON, string(contentAfter))
}
