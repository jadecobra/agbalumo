# Phase 6: Performance Audit Migration

## Objective
Port the logic from `scripts/performance-audit.sh` into a new `perf` Cobra subcommand in `cmd/verify/main.go`. This will deprecate the 428-line legacy bash script.

## Context
The performance audit performs file size checks (CSS, JS, images), static analysis via grep on Go source code for caching/SQLite pragmas, and live endpoint checks. Translating this to Go makes it cross-platform and heavily robust.

## Steps for Execution
1. Open `cmd/verify/main.go`.
2. Add a `perfCmd` Cobra command:
   - Use `os.Stat` to check file sizes of `ui/static/css/output.css` (< 150KB) and `ui/static/js/app.js` (< 50KB).
   - Use `os.ReadFile` and `strings.Contains` to check `internal/repository/sqlite/sqlite.go` for critical pragmas: `journal_mode=WAL`, `busy_timeout`, `MaxOpenConns`.
   - Execute the smoke benchmark natively: `runCmd("go", "test", "-bench=BenchmarkSearchPerformance/FindAll_Default_Page1", "-benchtime=100ms", "./internal/repository/sqlite/search_performance_test.go")`.
3. Register `perfCmd` in `rootCmd.AddCommand()`.
4. Run `go run cmd/verify/main.go perf` to ensure parity.
5. Create an artifact `performance_report.md` manually to log any findings.
6. Delete `scripts/performance-audit.sh`.
7. Commit changes natively: `refactor(ci): migrate performance audit out of bash into verify cli`.

## Verification
- `scripts/performance-audit.sh` should no longer exist.
- `go run cmd/verify/main.go perf` should correctly print the performance audit results.
