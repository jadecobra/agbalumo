# Task 23: Module-Level Cleanup (Listing)

## Context
The `internal/module/listing/` package is the largest contributor to test duplication, with over 100 clone groups identified in `_test.go` files (create, update, delete, search, etc.). This task performs a systematic cleanup.

## Checklist
- [ ] Consolidate the 20+ clones of `SaveListing` and `NewListing` calls into shared helpers in `listing_helpers_test.go`.
- [ ] Refactor `listing_create_test.go`, `listing_update_test.go`, and `listing_search_test.go` to use these helpers.
- [ ] Consolidate shared UI-related test assertions (e.g., checking for specific HTML tags) into `internal/testutil/ui.go`.

## Verification
- [ ] Run `go test ./internal/module/listing/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` and ensure significant reduction in the clone count for `listing` module.
