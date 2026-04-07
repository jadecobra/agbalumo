# Task 22: Module-Level Cleanup (Admin)

## Context
The `internal/module/admin/` package contains 20+ clones across actions, bulk, and dashboard tests. This task performs a comprehensive cleanup by moving shared setup and teardown logic to `internal/module/admin/test_helpers_test.go`.

## Checklist
- [ ] Move repetitive `NewAdminHandler`, `setupContext`, and `mockRepository` calls into shared helpers.
- [ ] Refactor `admin_actions_test.go`, `admin_bulk_test.go`, and `admin_dashboard_test.go` to use these helpers.

## Verification
- [ ] Run `go test ./internal/module/admin/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` and ensure significant reduction in the clone count for `admin` module.
