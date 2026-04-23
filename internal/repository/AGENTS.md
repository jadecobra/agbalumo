# Repository Layer Intelligence

This package is responsible for all data access. When working here, adhere to the following strict constraints:

# Repository Constraints
- Production DB: SQLite with WAL mode
- All public queries MUST filter on `last_verified_at` (Zombie Data rule)
- Spatial queries use bounding-box pre-filter + Haversine
- Test with file-backed DB for WAL behavior, not just :memory:

## SQL & SQLite Specifics
*   **Pragmas**: Always ensure `PRAGMA foreign_keys = ON;` and `PRAGMA journal_mode = WAL;` are configured in `NewSQLiteRepository`.
*   **No Raw Concatenation**: NEVER use `fmt.Sprintf` or string concatenation to build SQL queries. Use `?` placeholders for all parameters to prevent SQL injection.
*   **Transaction Wrappers**: Any operation that modifies multiple tables or a single table in multiple steps MUST be wrapped in a transaction (`repo.db.Begin()`).
*   **FTS Integrity**: If you modify the `listings` table, you must verify that the FTS5 triggers (trigram search) are not broken.

## Data Mapping
*   **Null Handling**: Use `sql.NullString`, `sql.NullInt64`, etc., when mapping columns that are nullable in the schema.
*   **Timestamp Precision**: Store and retrieve all timestamps in UTC (`time.UTC`).
*   **Error Mapping**: Always map implementation-specific errors (SQLITE_CONSTRAINT, etc.) to domain errors like `domain.ErrConflict` or `domain.ErrNotFound`.

## Architecture Decisions
*   **No Business Logic**: The repository should be a "dumb" data access layer. Do not implement complex business rules or multi-step validations here; that belongs in `internal/service`.
*   **Migrations**: All schema changes must be implemented as separate `.sql` files in `internal/repository/sqlite/migrations` to ensure idempotency and Agent readability.
