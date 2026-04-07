# Task 18: Test Helper Standardization (UI)

## Context
Pervasive duplication of UI-related test helpers has been found in `listing_helpers_test.go` and `auth/test_helpers_test.go`. This task consolidates them into `internal/testutil/ui.go`.

## Checklist
- [ ] Consolidate shared UI-related test helpers into `internal/testutil/ui.go`.
- [ ] Refactor `listing_helpers_test.go`, `auth/test_helpers_test.go`, and other UI-related tests to use these shared helpers.

## Verification
- [ ] Run `go test ./...` and ensure all tests pass.
- [ ] Verify that UI duplication has been significantly reduced across all modules.
