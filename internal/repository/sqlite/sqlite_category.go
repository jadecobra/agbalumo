package sqlite

import (
	"context"
	"database/sql"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// SaveCategory inserts or updates a category.
func (r *SQLiteRepository) SaveCategory(ctx context.Context, c domain.CategoryData) error {
	query := `
	INSERT INTO categories (id, name, claimable, is_system, active, requires_special_validation, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		claimable = excluded.claimable,
		is_system = excluded.is_system,
		active = excluded.active,
		requires_special_validation = excluded.requires_special_validation,
		updated_at = excluded.updated_at;
	`
	_, err := r.writeDB.ExecContext(ctx, query,
		c.ID, c.Name, c.Claimable, c.IsSystem, c.Active, c.RequiresSpecialValidation, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

// GetCategories retrieves categories based on the provided filter.
func (r *SQLiteRepository) GetCategories(ctx context.Context, filter domain.CategoryFilter) ([]domain.CategoryData, error) {
	query := `
		SELECT id, name, claimable, is_system, active, requires_special_validation, created_at, updated_at
		FROM categories
		WHERE 1=1
	`
	var args []interface{}

	if filter.ActiveOnly {
		query += ` AND active = 1`
	}

	query += ` ORDER BY name ASC`

	rows, err := r.readDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var categories []domain.CategoryData
	for rows.Next() {
		var c domain.CategoryData
		var created, updated sql.NullTime
		err := rows.Scan(&c.ID, &c.Name, &c.Claimable, &c.IsSystem, &c.Active, &c.RequiresSpecialValidation, &created, &updated)
		if err != nil {
			return nil, err
		}
		if created.Valid {
			c.CreatedAt = created.Time
		}
		if updated.Valid {
			c.UpdatedAt = updated.Time
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

// EnsureCoreCategories seeds the categories table with core types that must exist.
func (r *SQLiteRepository) UpsertCoreCategory(ctx context.Context, c domain.CategoryData) error {
	query := `
	INSERT INTO categories (id, name, claimable, is_system, active, requires_special_validation, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		claimable = excluded.claimable,
		is_system = excluded.is_system,
		active = excluded.active,
		requires_special_validation = excluded.requires_special_validation,
		updated_at = excluded.updated_at;
	`
	_, err := r.writeDB.ExecContext(ctx, query,
		c.ID, c.Name, c.Claimable, c.IsSystem, c.Active, c.RequiresSpecialValidation, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

// GetCategory retrieves a single category by its name (ID).
func (r *SQLiteRepository) GetCategory(ctx context.Context, name string) (domain.CategoryData, error) {
	query := `
		SELECT id, name, claimable, is_system, active, requires_special_validation, created_at, updated_at
		FROM categories
		WHERE id = ?
	`
	row := r.readDB.QueryRowContext(ctx, query, name)

	var c domain.CategoryData
	var created, updated sql.NullTime
	err := row.Scan(&c.ID, &c.Name, &c.Claimable, &c.IsSystem, &c.Active, &c.RequiresSpecialValidation, &created, &updated)
	if err == sql.ErrNoRows {
		return domain.CategoryData{}, domain.ErrCategoryNotFound
	}
	if err != nil {
		return domain.CategoryData{}, err
	}
	if created.Valid {
		c.CreatedAt = created.Time
	}
	if updated.Valid {
		c.UpdatedAt = updated.Time
	}
	return c, nil
}
