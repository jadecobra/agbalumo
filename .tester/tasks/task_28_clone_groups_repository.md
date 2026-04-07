# Task 28: Clone Group Reduction — Batch C (repository & SQLite test seeding)

## Context

The SQLite and repository test files are the single largest contributor to clone
groups by count. Most are caused by inline DB seeding (the same `newTestListing()` +
`repo.SaveListing(...)` pattern repeated across 15+ test files). This task consolidates
them.

---

## Target Clone Groups

### 1. `sqlite_listing_search_test.go` — 15-clone single-line seeding pattern (lines 22–79)

**Pattern**: `ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)`
(or similar) followed by `repo.SaveListing(ctx, newTestListing(...))` — repeated
for every sub-test.

**Fix**: Already have `newTestListing` helper? If not, add it to
`internal/repository/sqlite/sqlite_test_helpers_test.go`. Then replace all
inline saves with:
```go
saveTestListing(t, ctx, repo, domain.Listing{...})
```
where `saveTestListing` wraps the call and calls `t.Fatal` on error.

### 2. `sqlite_listing_bench_test.go:17–50` — 2-clone bench setup block

**Pattern**: Two benchmark functions repeat the same 11-line setup block
(open DB, insert N listings, defer cleanup).

**Fix**: Extract `setupBenchmarkDB(b *testing.B, n int) (Repository, func())` helper.
Both `BenchmarkX` and `BenchmarkY` call the helper.

### 3. `repro_test.go:30,55` and `sqlite_category_test.go:104–291` — 10-clone category save

**Pattern**: `repo.SaveCategory(ctx, domain.CategoryData{...})` repeated inline
without a helper.

**Fix**: Add `saveTestCategory(t, ctx, repo, cat)` to the test helper file
(already likely exists or was created in task 22). Use it throughout.

### 4. `seeder/category_seeder_test.go:46,97` and `seeder/config_verification_test.go:34,57` — 10-clone overlap

**Pattern**: Seeder tests repeat the same `NewCategorySeeder(...)` + `Seed(ctx)`
setup block.

**Fix**: Add `setupSeeder(t *testing.T) (*CategorySeeder, func())` test helper
in `internal/seeder/seeder_test_helpers_test.go`.

### 5. `internal/repository/cached/cached_test.go:106–194` — 2-clone cache assertion block

**Pattern**: `assert.Equal(t, expected, got)` + `assert.Equal(t, 1, spy.calls)` 
repeated for two scenarios (hit vs miss).

**Fix**: Extract `assertCacheResult(t, spy, result, wantCalls int)` helper in the test file.

### 6. `internal/repository/sqlite/sqlite_stats.go:53–72` — 2-clone stat query block

**Pattern**: Two near-identical SQL query-and-scan sequences for daily metrics.

**Fix**: Extract `queryDailyMetrics(ctx, db, query string) ([]domain.DailyMetric, error)`
private function. Both `GetListingGrowth` and `GetUserGrowth` call it with different SQL.

---

## Verification

- [ ] `go test ./internal/repository/... ./internal/seeder/...`
- [ ] `go run cmd/verify/main.go critique` — confirm repository clone count is reduced
- [ ] Commit: `refactor(repository): consolidate test seeding helpers to reduce clone groups`
