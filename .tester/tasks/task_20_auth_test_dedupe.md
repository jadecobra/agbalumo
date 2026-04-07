# Task 20: Auth Handler Test Consolidation

## Context
The `internal/module/auth/` test files (especially `handler_register_test.go`) contain 15+ nearly identical registration request setups and response assertions. This task consolidates them into shared test helpers.

## Checklist
- [ ] Identify repetitive `http.NewRequest` and `httptest.NewRecorder` blocks in `handler_register_test.go`.
- [ ] Create a `performRegistration(t *testing.T, payload map[string]string) *httptest.ResponseRecorder` helper.
- [ ] Refactor all registration tests to use this helper.

## Verification
- [ ] Run `go test ./internal/module/auth/...` and ensure all tests pass.
- [ ] Run `go run cmd/verify/main.go critique` to confirm the clone group count decreased for auth.
