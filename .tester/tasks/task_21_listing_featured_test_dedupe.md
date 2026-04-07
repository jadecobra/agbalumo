# Task 21: Listing Featured Test Consolidation

## Context
The `internal/module/listing/listing_featured_test.go` file contains 40+ clones of featured status assertions. This task refactors these assertions into a single `assertFeaturedStatus(t *testing.T, id string, expected bool)` helper.

## Checklist
- [ ] Implement `assertFeaturedStatus(t *testing.T, id string, expected bool)` in `listing_featured_test.go`.
- [ ] Refactor all featured status checks to use this helper.

## Verification
- [ ] Run `go test ./internal/module/listing/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` and ensure significant reduction in the clone count for `listing_featured_test.go`.
