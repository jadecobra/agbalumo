package testutil

import (
	"io"
	"log/slog"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

// SetupTestRepository initializes an in-memory sqlite repository for integration tests.
func SetupTestRepository(t *testing.T) *sqlite.SQLiteRepository {
	t.Helper()
	// Using ":memory:" creates a new, empty, in-memory database every time.
	// Since we set MaxOpenConns(1) in the repository, this works well for a single instance.
	repo, err := sqlite.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory sqlite repository: %v", err)
	}
	return repo
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
	)

	return app, func() {
		_ = repo.Close()
	}
}
