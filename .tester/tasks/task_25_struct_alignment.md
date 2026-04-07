# Task 25: Struct Alignment (Gate Failure)

## Context

`go run cmd/verify/main.go critique` [3/4] fails because 28 structs across the
codebase are sub-optimally field-ordered, wasting between 8 and 48 bytes per
instance due to alignment padding.

The fix is **purely mechanical**: reorder struct fields from largest to smallest
alignment (pointers/interfaces → int64/float64 → int32/float32 → int16 → int8/bool).
No logic changes. No test changes needed.

**Tool**: `go run golang.org/x/tools/cmd/fieldalignment -fix ./...`

> ⚠️ WARNING: `fieldalignment -fix` rewrites files in-place. Always run tests after.

---

## Checklist

- [ ] **Run `fieldalignment` auto-fix across the whole repo**
  ```
  go run golang.org/x/tools/cmd/fieldalignment -fix ./...
  ```
  This will reorder fields in all flagged structs automatically.

- [ ] **Verify no test regressions** (struct field reordering cannot change behavior,
  but literal struct initializers that use positional (non-named) fields will break
  the build — check for those first):
  ```
  grep -rn "domain\.CategoryData{[^}]*}" . --include="*.go" | grep -v "Name:"
  ```
  If any positional initializers are found, update them to use named fields before
  running the fix.

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go run cmd/verify/main.go critique` — confirm `[3/4] Struct Alignment` now shows ✅

---

## Key Files (28 total)

**Production** (highest savings):
| File | Savings |
|------|---------|
| `internal/maintenance/cost.go:22` | 56 → 8 bytes |
| `internal/domain/category.go:23` | 88 → 72 bytes |
| `internal/domain/category.go:40` | 72 → 32 bytes |
| `internal/repository/cached/cached.go:12` | 120 → 80 bytes |
| `internal/maintenance/watcher.go:16` | 48 → 32 bytes |
| `internal/domain/claim.go:15` | 136 → 128 bytes |
| `internal/domain/csv.go:9` | 32 → 8 bytes |
| `internal/domain/feedback.go:13` | 88 → 80 bytes |
| `internal/domain/listing.go:24` | 24 → 16 bytes |
| `internal/config/config.go:10` | 120 → 104 bytes |
| `internal/repository/sqlite/sqlite.go:20` | 24 → 16 bytes |
| `cmd/serve.go:19` | 48 → 40 bytes |

**Test files** (also flagged — fix with same command):
`listing_validation_basic_test.go`, `listing_validation_extended_test.go`,
`admin_actions_test.go`, `admin_login_test.go`, `pagination_test.go`,
`listing_create_test.go`, `listing_delete_test.go`, `listing_form_integration_test.go`,
`listing_update_test.go` (2 structs), `cached_test.go`, `geocoding_test.go`,
`seed_test.go`, `serve_test.go`, `server_public_test.go`

---

## Verification

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go run cmd/verify/main.go critique` — `[3/4]` must show ✅
- [ ] Commit: `refactor(domain): fix struct field alignment for memory efficiency`
