# Task 16: Domain Validation Consolidation

## Context
The `internal/domain/listing.go` file contains 5+ clones of validation rules (e.g., checking if title, description, or id is empty). This task refactors these into a reusable rules engine/map to reduce line count and complexity.

## Checklist
- [ ] Create a `validationRules` slice or map in `internal/domain/listing.go`.
- [ ] Refactor `Validate()` to loop through these rules.
- [ ] Consolidate redundant string checks.

## Verification
- [ ] Run `go test ./internal/domain/...` and ensure all validation tests pass.
- [ ] Run `go run cmd/verify/main.go critique` and ensure the clone group count has decreased for domain.
