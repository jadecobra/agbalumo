# Agents: internal/repository/sqlite

## OVERVIEW

SQLite persistence layer implementing domain store interfaces with auto-migration.

## WHERE TO LOOK

- `sqlite.go` — Main repository: Listings, Users, Feedback CRUD, metrics queries
- `feedback.go` — Feedback-specific operations (extends repository)
- `*_test.go` — Co-located tests with table-driven test cases

## CONVENTIONS

```go
// Constructor pattern
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error)
func NewSQLiteRepositoryFromDB(db *sql.DB) *SQLiteRepository

// All methods accept context.Context as first parameter
func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error

// Nullable fields use sql.NullTime, sql.NullString
var deadline sql.NullTime
if deadline.Valid { l.Deadline = deadline.Time }

// Queries use ExecContext/QueryContext/QueryRowContext
rows, err := r.db.QueryContext(ctx, query, args...)

// Scan to domain types via helper functions
func scanListing(s Scanner) (domain.Listing, error)
func scanUser(s Scanner) (domain.User, error)

// Migration runs on startup via ALTER TABLE (add columns if missing)
// PRAGMA tuning: WAL mode, busy_timeout=5000, synchronous=NORMAL
```

## ANTI-PATTERNS

- Do NOT use `database/sql` directly in handlers — go through repository
- Do NOT hardcode connection strings — use config/dbpath
- Do NOT skip `defer rows.Close()` after QueryContext
- Do NOT ignore `rows.Err()` after iteration loop
- Do NOT store time as strings — use DATETIME columns with Go time.Time
