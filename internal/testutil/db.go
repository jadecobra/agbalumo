package testutil

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"

	"context"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"time"
)

var (
	dbCounter int64
	counterMu sync.Mutex
)

// SetupTestRepository initializes an in-memory sqlite repository for integration tests.
func SetupTestRepository(t testing.TB) *sqlite.SQLiteRepository {
	// Using ":memory:" creates a new, empty, in-memory database every time.
	// Since we set MaxOpenConns(1) in the repository, this works well for a single instance.
	repo, err := sqlite.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory sqlite repository: %v", err)
	}
	return repo
}

// SetupTestDB returns a raw, unmigrated in-memory sqlite connection.
func SetupTestDB(t testing.TB) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?_time_format=sqlite")
	if err != nil {
		t.Fatalf("failed to open unmigrated in-memory db: %v", err)
	}
	return db
}

// SetupTestRepositoryUnique initializes an in-memory sqlite repository with a unique name.
// This ensures isolation between parallel tests while allowing the repository to use
// separate read and write connection pools even in-memory (via shared cache mode).
func SetupTestRepositoryUnique(t testing.TB) (*sqlite.SQLiteRepository, string) {
	counterMu.Lock()
	dbCounter++
	id := dbCounter
	counterMu.Unlock()

	// Use a unique name per test to ensure isolation in shared-cache mode.
	// We include _time_format=sqlite to ensure consistent date/time parsing.
	dbName := fmt.Sprintf("file:test_%s_%d?mode=memory&cache=shared&_time_format=sqlite", t.Name(), id)
	repo, err := sqlite.NewSQLiteRepository(dbName)
	if err != nil {
		t.Fatalf("failed to create unique in-memory repository: %v", err)
	}
	return repo, dbName
}

// SetupTestAppEnv initializes a test AppEnv with an in-memory database and functional mocks.
func SetupTestAppEnv(t *testing.T) (*env.AppEnv, func()) {
	t.Helper()
	repo, err := sqlite.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory sqlite repository: %v", err)
	}

	cfg := config.LoadConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app := env.NewAppEnv(
		repo,
		cfg,
		logger,
		&MockCSVService{},
		&MockGeocodingService{},
		&StubImageService{},
		&MockListingService{},
		&MockCategorizationService{},
		&MockMetricsService{},
	)

	return app, func() {
		_ = repo.Close()
	}
}

// SaveTestListing creates and saves a listing with the given identifier and title.
func SaveTestListing(t testing.TB, db domain.ListingRepository, id, title string, extra ...func(*domain.Listing)) {
	t.Helper()
	l := domain.Listing{
		ID:           id,
		Title:        title,
		Type:         domain.Business,
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
		CreatedAt:    time.Now(),
	}
	for _, f := range extra {
		f(&l)
	}
	if err := db.Save(context.Background(), l); err != nil {
		t.Fatalf("Failed to save test listing %s: %v", id, err)
	}
}

// AssertFeaturedStatus verifies the featured status of a listing in the database.
func AssertFeaturedStatus(t testing.TB, db domain.ListingRepository, id string, expected bool) {
	t.Helper()
	l, err := db.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("listing %s retrieval failed: %v", id, err)
	}
	if l.Featured != expected {
		t.Errorf("listing %s featured status mismatch: expected %v, got %v", id, expected, l.Featured)
	}
}
