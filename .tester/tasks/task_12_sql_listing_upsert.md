# Task 12: SQL Consolidation (Listing UPSERT)

## Context
The `ChiefCritic` audit identified multiple clones of the `listings` UPSERT SQL query in `internal/repository/sqlite/sqlite_listing_write.go`. This task centralizes that logic into a package-local `queries.go` file.

## Checklist
- [ ] Create `internal/repository/sqlite/queries.go`.
- [ ] Define `const ListingUpsertSQL = ...` in `queries.go`.
- [ ] Refactor `internal/repository/sqlite/sqlite_listing_write.go` to use this constant.

## Verification
- [ ] Run `go test ./internal/repository/sqlite/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` to confirm the number of clone groups decreased.
