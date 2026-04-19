# SQLite Repository: Agent Guidance

This package implements the `domain.ListingStore` and `domain.UserStore` interfaces using SQLite3.

## Architectural Constraints

- **Single Writer**: SQLite supports multiple readers but only one writer. Ensure transactions are used where atomicity is required.
- **FTS5 Integration**: All search-related logic must leverage the `listing_fts` virtual table. Do NOT implement manual string filtering in Go if FTS can handle it.
- **Parameter Binding**: NEVER use string concatenation for SQL queries. Always use `?` placeholders.
- **Bulk Operations**: Use the `BulkInsertListings` method for many-to-one insertions to minimize transaction overhead.

## Testing Patterns

- **In-Memory for Tests**: Tests should default to `:memory:` or a temporary file that is cleaned up.
- **Golden File Verification**: Use golden files for complex query results to ensure stability.

## Common Pitfalls

- **Constraint Violations**: Catch `SQLITE_CONSTRAINT` errors and map them to `domain.ErrConflict` or similar.
- **Triggers**: Be aware of the `listing_au` and `listing_ad` triggers that keep the FTS index in sync.
