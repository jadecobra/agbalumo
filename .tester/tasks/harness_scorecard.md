# Agent Harness Scorecard

> Weekly review of harness quality. Score each dimension 1-10.
> Run `./scripts/agent-exec.sh cost --text` and `go test ./internal/agent/... -cover` before scoring.

## Scoring Criteria

| Dimension | What to evaluate |
|---|---|
| **Architecture** | Clean separation, no circular calls, minimal coupling |
| **Security** | Anti-cheat effectiveness, signature integrity, bypass resistance |
| **Testability** | Error paths testable, no `os.Exit` in handlers, mock-friendly |
| **Maintainability** | DRY compliance, typed structs, no magic strings |
| **Completeness** | All gates functional, drift detection accurate, progress tracking works |
| **Dogfooding** | Harness meets its own coverage/lint standards |

## Weekly Scores

<!-- Copy the template row, fill in scores, and note what changed that week. -->

| Week | Date | Arch | Sec | Test | Maint | Compl | Dogfood | **Avg** | Notes |
|---|---|---|---|---|---|---|---|---|---|
| 0 (Baseline) | 2026-03-24 | 7 | 5 | 4 | 6 | 8 | 3 | **5.5** | Initial review. See `harness_review.md` for details. |
| 1 | | | | | | | | | |
| 2 | | | | | | | | | |
| 3 | | | | | | | | | |
| 4 | | | | | | | | | |

## Review Checklist

Before filling in the weekly row:

- [ ] Run `go test ./internal/agent/... -cover` — note coverage %
- [ ] Run `go test ./cmd/harness/... -cover` — note coverage %
- [ ] Run `./scripts/agent-exec.sh cost --text` — note RMS trend
- [ ] Check `progress.json` "Agent Harness Tech Debt" — count completed vs pending
- [ ] Spot-check: attempt `./scripts/agent-exec.sh gate red-test PASS` — should be blocked
- [ ] Review any new files in `internal/agent/` or `cmd/harness/commands/` for standards compliance
