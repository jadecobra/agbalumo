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

func TestMigration_OriginPriorityAndIndex(t *testing.T) {
	t.Parallel()
	repo, err := NewSQLiteRepository(":memory:")
	assert.NoError(t, err)
	defer func() { _ = repo.Close() }()

	var colExists int
	err = repo.writeDB.QueryRow("SELECT COUNT(*) FROM pragma_table_xinfo('listings') WHERE name = 'origin_priority'").Scan(&colExists)
	assert.NoError(t, err)
	assert.Equal(t, 1, colExists, "Column 'origin_priority' should exist in 'listings' table")

	var idxExists int
	err = repo.writeDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_listings_default_feed'").Scan(&idxExists)
	assert.NoError(t, err)
	assert.Equal(t, 1, idxExists, "Index 'idx_listings_default_feed' should exist")

	// Verify the default feed query plan uses idx_listings_default_feed and no temp B-tree
	query := `EXPLAIN QUERY PLAN SELECT rowid FROM listings WHERE 1=1 AND is_active = 1 AND status = 'Approved' ORDER BY featured DESC, origin_priority ASC, heat_level DESC, rating DESC, created_at DESC, rowid ASC LIMIT 30 OFFSET 0`
	rows, err := repo.readDB.Query(query)
	assert.NoError(t, err)
	defer func() { _ = rows.Close() }()

	var planLines []string
	for rows.Next() {
		var id, parent, notused int
		var detail string
		err = rows.Scan(&id, &parent, &notused, &detail)
		assert.NoError(t, err)
		planLines = append(planLines, detail)
	}
	assert.NoError(t, rows.Err())
	t.Logf("Default feed query plan: %v", planLines)

	for _, line := range planLines {
		assert.NotContains(t, line, "USE TEMP B-TREE FOR ORDER BY", "Default feed query should not use temporary B-tree for sorting")
	}
}

func TestParseNullableTime_Formats(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"2026-09-23T15:04:05Z", true},
		{"2026-09-23T15:04:05.123456789Z", true},
		{"2026-09-23 15:04:05", true},
		{"2026-09-23 15:04:05.999999999 -0700 MST", true},
		{"2026-09-23 15:04:05.999999999 -0700 MST m=+0.001", true},
		{"2026-09-23 15:04:05-07:00", true},
		{"invalid-time-format", false},
	}

	for _, tc := range cases {
		got := parseNullableTime(tc.input)
		if tc.want && got == nil {
			t.Errorf("parseNullableTime(%q) expected non-nil, got nil", tc.input)
		}
		if !tc.want && got != nil {
			t.Errorf("parseNullableTime(%q) expected nil, got %v", tc.input, got)
		}
	}
}
