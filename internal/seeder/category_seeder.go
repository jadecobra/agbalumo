package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

// EnsureCategoriesSeeded reads config/categories.json and upserts into the database.
func EnsureCategoriesSeeded(ctx context.Context, repo *sqlite.SQLiteRepository, configPath string) error {
	slog.Info("Starting category seed from config", "path", configPath)
	
	// #nosec G304 - Config path is controlled by application startup logic
	data, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("Categories config file not found, skipping category seed", "path", configPath)
			return nil
		}
		return fmt.Errorf("failed to read categories config: %w", err)
	}

	var categories []domain.CategoryData
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("failed to parse categories config: %w", err)
	}

	for _, c := range categories {
		err := repo.UpsertCoreCategory(ctx, c)
		if err != nil {
			slog.Error("Failed to upsert core category", "id", c.ID, "error", err)
		} else {
			slog.Debug("Upserted core category", "name", c.Name)
		}
	}

	slog.Info("Completed category seed")
	return nil
}
