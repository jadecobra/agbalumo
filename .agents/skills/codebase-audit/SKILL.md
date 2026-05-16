---
name: "Codebase Audit"
description: "Systematic 6-dimension health check — deterministic Phase 1 aggregating all verify tools."
triggers:
  - "audit codebase"
  - "health check"
  - "score the codebase"
  - "review infrastructure"
  - "how healthy is the codebase"
mutating: false
---

# Codebase Audit Skill

## Phase 1 — Deterministic Scan (Zero Reasoning Cost)

Run every verification gate and record pass/fail:

Run `go run ./cmd/verify sweep --json > /tmp/sweep.json`. Use the JSON output for Phase 2 scoring.

Collect results into a table:
```
| Gate | Result | Details |
|------|--------|---------|
```

## Phase 2 — Scoring (Minimal Reasoning)

Map gate results to 6 dimensions:

| Dimension | Gates | Score Rule |
|-----------|-------|-----------|
| Discovery & Navigation | agents-coverage, resolve | 10 if both pass, -1 per failure |
| Verification Tooling | skill-conformance, check-resolvable | 10 if both pass |
| Documentation Integrity | doc-drift | 10 if pass, count violations for partial |
| Architecture & Code | deprecated, template-contract | 10 if both pass, -1 per 5 violations |
| Skills & Workflows | skill-conformance, check-resolvable | 10 if both pass |
| Quota Efficiency | context-cost | 10 if RMS < 300, 8 if < 500, 6 if > 500 |

## Phase 3 — Recommendations (Reasoning)

Only if Phase 2 overall < 9.0:
1. Identify the lowest-scoring dimension
2. Propose ≤3 concrete fixes (grep-anchored, not vague)
3. Estimate token savings per fix

## Output
Produce a scored artifact with the table and any recommendations.
