# Task 17: Service Layer Extraction (Categorization)

## Context
The `admin` and `listing` modules both contain logic to fetch categories or count listings. This duplication violates our DDD hexagonal architecture. This task moves that logic into the `service` layer.

## Checklist
- [ ] Create/Update `internal/service/categorization.go` (if not already there).
- [ ] Implement `GetActiveCategories(ctx context.Context)` in the service layer.
- [ ] Refactor `internal/module/admin/admin.go` and `internal/module/listing/listing.go` to use this service.

## Verification
- [ ] Run `go test ./internal/module/admin/... ./internal/module/listing/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go ci` and confirm no regressions.
