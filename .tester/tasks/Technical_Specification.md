# Technical Review: Context Cost Optimization
**Persona**: @SystemsArchitect
**Feature**: `context-cost-optimization`

## 1. Context Cost Audit Analysis
The current RMS of **192.26** is heavily skewed by `critique_report.md` (2651 lines). Excluding this file alone will likely drop the RMS by 40-50%.

## 2. Refactoring Proposals

### A. Modularize Security Checks (`internal/agent/security.go`)
The current file is 628 lines and mixes raw regex scanning with complex AST traversal for Go code.
**Proposed Split**:
- `security.go`: Core types (`SecurityViolation`), public API (`VerifySecurityStatic`), and dispatcher (`checkFile`).
- `security_regex.go`: Global pattern definitions and raw string scanners (`checkSecretsRaw`, `checkStructuralRaw`).
- `security_ast.go`: All AST-based logic (`checkSQLi`, `checkXSS`, `checkSSRF`, etc.) and helper `isUnsafeString`.

### B. Split Large Test Files
The test files `verify_apispec_test.go` (641 lines) and `sqlite_listing_test.go` (578 lines) have reached a point where they contribute significantly to context noise.
**Proposed Split**:
- `internal/agent/verify_apispec_test.go` -> `internal/agent/verify_apispec_routes_test.go` & `internal/agent/verify_apispec_infra_test.go`.
- `internal/repository/sqlite/sqlite_listing_test.go` -> `internal/repository/sqlite/listing_crud_test.go` & `internal/repository/sqlite/listing_search_test.go`.

### C. Cost Tool Logic Enhancement
- Update `internal/agent/cost.go` to support a glob-based `ignoredPatterns` list.
- Add `*.report.md` and `*.log` to the default ignore set.

## 3. Risk Assessment
- **Circular Dependencies**: Splitting `security.go` within the same package is safe but requires careful placement of `var` blocks.
- **Test Integrity**: Moving tests must ensure that `go test ./...` still picks up all cases. Using the `_test.go` suffix is mandatory.

## 4. Verification Strategy
1. Run `./scripts/agent-exec.sh cost` before and after each major refactor.
2. Target RMS: **< 110**.
3. All existing security and contract tests must pass.
