package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestSaveAndFindByID(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	l := domain.Listing{
		ID:           "test-1",
		Title:        "Original Title",
		OwnerOrigin:  "Ghana",
		Type:         domain.Business,
		IsActive:     true,
		CreatedAt:    time.Now(),
		ContactEmail: "test@example.com",
	}

	ctx := context.Background()

	// 1. Save New
	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// 2. Find
	found, err := repo.FindByID(ctx, "test-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
	if found.Title != "Original Title" {
		t.Errorf("Expected title 'Original Title', got '%s'", found.Title)
	}

	// 3. Update (Save existing)
	l.Title = "Updated Title"
	err = repo.Save(ctx, l)
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	foundUpdated, err := repo.FindByID(ctx, "test-1")
	if err != nil {
		t.Fatalf("Failed to find updated: %v", err)
	}
	if foundUpdated.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", foundUpdated.Title)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()
	_ = repo.Save(ctx, domain.Listing{ID: "del-me", Title: "Delete Me"})

	err := repo.Delete(ctx, "del-me")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(ctx, "del-me")
	if err == nil {
		t.Error("Expected error finding deleted listing, got nil")
	}

	// Delete non-existent
	err = repo.Delete(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error deleting non-existent listing, got nil")
	}
}

func TestHoursOfOperationPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	l := domain.Listing{
		ID:               "h-1",
		Title:            "Hours Test",
		HoursOfOperation: "Mon-Fri 9am-5pm",
		IsActive:         true,
	}

	_ = repo.Save(ctx, l)

	found, _ := repo.FindByID(ctx, "h-1")
	if found.HoursOfOperation != l.HoursOfOperation {
		t.Errorf("Expected hours %q, got %q", l.HoursOfOperation, found.HoursOfOperation)
	}
}

func TestCategoryErrors_Raw(t *testing.T) {
	t.Parallel()
	db, _ := sql.Open("sqlite", ":memory:")
	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	_ = db.Close()
	ctx := context.Background()

	_, err := repo.GetCategory(ctx, "any")
	if err == nil {
		t.Error("Expected error on closed DB")
	}
}

func TestEnrichmentAttemptedAtPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	now := time.Now().Truncate(time.Second) // SQLite might truncate sub-second precision
	l := domain.Listing{
		ID:                    "enrich-1",
		Title:                 "Enrich Test",
		EnrichmentAttemptedAt: &now,
		IsActive:              true,
	}

	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "enrich-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}

	if found.EnrichmentAttemptedAt == nil {
		t.Fatal("Expected EnrichmentAttemptedAt to be set, got nil")
	}

	if !found.EnrichmentAttemptedAt.Equal(now) {
		t.Errorf("Expected EnrichmentAttemptedAt %v, got %v", now, *found.EnrichmentAttemptedAt)
	}
}

func TestFindEnrichmentTargets_FiltersAttemptedAt(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)
	eightDaysAgo := now.AddDate(0, 0, -8)

	l1 := domain.Listing{
		ID:         "unenriched-1",
		Title:      "Unenriched 1",
		WebsiteURL: "http://test.com",
		IsActive:   true,
	}
	l2 := domain.Listing{
		ID:                    "unenriched-2",
		Title:                 "Unenriched 2",
		WebsiteURL:            "http://test.com",
		EnrichmentAttemptedAt: &oneDayAgo,
		IsActive:              true,
	}
	l3 := domain.Listing{
		ID:                    "unenriched-3",
		Title:                 "Unenriched 3",
		WebsiteURL:            "http://test.com",
		EnrichmentAttemptedAt: &eightDaysAgo,
		IsActive:              true,
	}

	_ = repo.Save(ctx, l1)
	_ = repo.Save(ctx, l2)
	_ = repo.Save(ctx, l3)

	targets, err := repo.FindEnrichmentTargets(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to find enrichment targets: %v", err)
	}

	targetsMap := make(map[string]bool)
	for _, target := range targets {
		targetsMap[target.ID] = true
	}

	if !targetsMap["unenriched-1"] {
		t.Error("Expected to find unenriched-1 (attempted_at is NULL)")
	}
	if targetsMap["unenriched-2"] {
		t.Error("Did NOT expect to find unenriched-2 (attempted_at is 1 day ago)")
	}
	if !targetsMap["unenriched-3"] {
		t.Error("Expected to find unenriched-3 (attempted_at is 8 days ago)")
	}
}

func TestDeliveryPlatformsPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	l := domain.Listing{
		ID:                "dp-1",
		Title:             "Delivery Test",
		DeliveryPlatforms: `["UberEats", "DoorDash"]`,
		IsActive:          true,
	}

	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "dp-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
	if found.DeliveryPlatforms != l.DeliveryPlatforms {
		t.Errorf("Expected delivery platforms %q, got %q", l.DeliveryPlatforms, found.DeliveryPlatforms)
	}
}

func TestQualityProxyPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	now := time.Now().Truncate(time.Second)
	l := domain.Listing{
		ID:              "qp-1",
		Title:           "Quality Test",
		Rating:          4.6,
		ReviewCount:     120,
		RatingUpdatedAt: &now,
		IsActive:        true,
	}

	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "qp-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
	if found.Rating != l.Rating {
		t.Errorf("Expected rating %f, got %f", l.Rating, found.Rating)
	}
	if found.ReviewCount != l.ReviewCount {
		t.Errorf("Expected review count %d, got %d", l.ReviewCount, found.ReviewCount)
	}
	if found.RatingUpdatedAt == nil {
		t.Fatal("Expected RatingUpdatedAt to be set, got nil")
	}
	if !found.RatingUpdatedAt.Equal(now) {
		t.Errorf("Expected RatingUpdatedAt %v, got %v", now, *found.RatingUpdatedAt)
	}
}

func TestDefaultSortingWithRating(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	l1 := domain.Listing{ID: "sort-1", Title: "Low Rating", HeatLevel: 5, Rating: 2.0, ReviewCount: 10, IsActive: true, Status: domain.ListingStatusApproved, Type: domain.Food}
	l2 := domain.Listing{ID: "sort-2", Title: "High Rating", HeatLevel: 5, Rating: 4.8, ReviewCount: 10, IsActive: true, Status: domain.ListingStatusApproved, Type: domain.Food}
	l3 := domain.Listing{ID: "sort-3", Title: "Mid Rating", HeatLevel: 5, Rating: 3.5, ReviewCount: 10, IsActive: true, Status: domain.ListingStatusApproved, Type: domain.Food}

	_ = repo.Save(ctx, l1)
	_ = repo.Save(ctx, l2)
	_ = repo.Save(ctx, l3)

	// Note: FindAll defaults to Category Food if empty, so we set Type to Food above.
	listings, _, err := repo.FindAll(ctx, string(domain.Food), "", "", 0, 0, 0, "", "", false, 10, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(listings) < 3 {
		t.Fatalf("Expected at least 3 listings, got %d", len(listings))
	}

	if listings[0].ID != "sort-2" {
		t.Errorf("Expected first listing to be sort-2, got %s", listings[0].ID)
	}
	if listings[1].ID != "sort-3" {
		t.Errorf("Expected second listing to be sort-3, got %s", listings[1].ID)
	}
	if listings[2].ID != "sort-1" {
		t.Errorf("Expected third listing to be sort-1, got %s", listings[2].ID)
	}
}

func TestGetLocationsWithCoordinates(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	l1 := domain.Listing{
		ID:        "loc-1",
		Title:     "Dallas Spot",
		City:      "Dallas",
		State:     "TX",
		Country:   "USA",
		Latitude:  32.7767,
		Longitude: -96.7970,
		IsActive:  true,
		Status:    domain.ListingStatusApproved,
		Type:      domain.Food,
	}
	l2 := domain.Listing{
		ID:        "loc-2",
		Title:     "Houston Spot",
		City:      "Houston",
		State:     "TX",
		Country:   "USA",
		Latitude:  29.7604,
		Longitude: -95.3698,
		IsActive:  true,
		Status:    domain.ListingStatusApproved,
		Type:      domain.Food,
	}

	if err := repo.Save(ctx, l1); err != nil {
		t.Fatalf("Failed to save l1: %v", err)
	}
	if err := repo.Save(ctx, l2); err != nil {
		t.Fatalf("Failed to save l2: %v", err)
	}

	locations, err := repo.GetLocations(ctx)
	if err != nil {
		t.Fatalf("GetLocations failed: %v", err)
	}

	assertLocationCoords(t, locations, "Dallas", 32.7767, -96.7970)
	assertLocationCoords(t, locations, "Houston", 29.7604, -95.3698)
}

func assertLocationCoords(t *testing.T, locations []domain.Location, city string, lat, lng float64) {
	t.Helper()
	loc, found := findCity(locations, city)
	if !found {
		t.Errorf("%s not found in locations list", city)
		return
	}
	if loc.Latitude == 0.0 || loc.Longitude == 0.0 {
		t.Errorf("Expected %s location to have non-zero coordinates", city)
	}
	if loc.Latitude != lat || loc.Longitude != lng {
		t.Errorf("Expected %s coords (%f, %f), got (%f, %f)", city, lat, lng, loc.Latitude, loc.Longitude)
	}
}

func findCity(locations []domain.Location, city string) (domain.Location, bool) {
	for _, loc := range locations {
		if loc.City == city {
			return loc, true
		}
	}
	return domain.Location{}, false
}
