# Task 14: SQL Consolidation (Read Queries)

## Context
The `sqlite_listing_read.go` and `sqlite_user.go` files contain several duplicated read and search queries. This task moves them into the centralized `queries.go` for consistency.

## Checklist
- [ ] Move common read queries for Listing, Category, and User into `internal/repository/sqlite/queries.go`.
- [ ] Refactor `internal/repository/sqlite/sqlite_listing_read.go` and `internal/repository/sqlite/sqlite_user.go` to use these constants.

## Verification
- [ ] Run `go test ./internal/repository/sqlite/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` to confirm the number of clone groups decreased.
