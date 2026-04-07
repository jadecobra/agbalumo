# Task 19: Test Helper Standardization (DB)

## Context
In-memory SQLite database setup logic is duplicated in several test files. This task moves that logic into `internal/testutil/db.go`.

## Checklist
- [ ] Move in-memory SQLite setup logic to `internal/testutil/db.go`.
- [ ] Refactor tests to use the shared database setup helper.

## Verification
- [ ] Run `go test ./...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` and confirm that test-related clone count has decreased.
