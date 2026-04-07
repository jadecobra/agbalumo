# Task 13: SQL Consolidation (Category/Claim UPSERT)

## Context
The `sqlite_category.go` and `sqlite_claim.go` files contain identical or near-identical UPSERT SQL blocks. This task moves these to the centralized `queries.go`.

## Checklist
- [ ] Move `CategoryUpsertSQL` and `ClaimUpsertSQL` from their respective write files to `internal/repository/sqlite/queries.go`.
- [ ] Refactor `internal/repository/sqlite/sqlite_category.go` and `internal/repository/sqlite/sqlite_claim.go` to use these constants.

## Verification
- [ ] Run `go test ./internal/repository/sqlite/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` to confirm the number of clone groups decreased.
