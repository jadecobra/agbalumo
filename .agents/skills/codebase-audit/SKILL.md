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

```bash
# Discovery & Navigation
go run ./cmd/verify agents-coverage        # AGENTS.md coverage %
go run ./cmd/verify resolve "test"          # Resolver working

# Verification Tooling
go run ./cmd/verify skill-conformance      # Skill YAML valid
go run ./cmd/verify check-resolvable       # All skills resolvable

# Documentation Integrity
go run ./cmd/verify doc-drift              # No stale references

# Architecture & Code
go run ./cmd/verify deprecated             # No deprecated patterns
go run ./cmd/verify template-contract      # ViewModel contracts valid

# Quota Efficiency
go run ./cmd/verify context-cost           # Token density
```

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
