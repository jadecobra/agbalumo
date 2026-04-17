# Repository Layer Intelligence

This package is responsible for all data access. When working here, adhere to the following strict constraints:

## SQL & SQLite Specifics
*   **Pragmas**: Always ensure `PRAGMA foreign_keys = ON;` and `PRAGMA journal_mode = WAL;` are configured in `NewSQLiteRepository`.
*   **No Raw Concatenation**: NEVER use `fmt.Sprintf` or string concatenation to build SQL queries. Use `?` placeholders for all parameters to prevent SQL injection.
*   **Transaction Wrappers**: Any operation that modifies multiple tables or a single table in multiple steps MUST be wrapped in a transaction (`repo.db.Begin()`).
*   **FTS Integrity**: If you modify the `listings` table, you must verify that the FTS5 triggers (trigram search) are not broken.

## Data Mapping
*   **Null Handling**: Use `sql.NullString`, `sql.NullInt64`, etc., when mapping columns that are nullable in the schema.
*   **Timestamp Precision**: Store and retrieve all timestamps in UTC (`time.UTC`).
