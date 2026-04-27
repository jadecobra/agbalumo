package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration_EnrichmentAttemptedAt(t *testing.T) {
	t.Parallel()
	repo, err := NewSQLiteRepository(":memory:")
	assert.NoError(t, err)
	defer func() { _ = repo.Close() }()

	var exists int
	err = repo.writeDB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('listings') WHERE name = 'enrichment_attempted_at'").Scan(&exists)
	assert.NoError(t, err)
	assert.Equal(t, 1, exists, "Column 'enrichment_attempted_at' should exist in 'listings' table")
}
