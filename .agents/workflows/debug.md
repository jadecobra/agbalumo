---
description: Structured diagnosis for CI failures, visual regressions, and environment drift.
---
# /debug <symptom>

## Phase 1: Scope (Deterministic)
1. `git diff --name-only HEAD~3 HEAD` — what changed recently?
2. `go run ./cmd/verify preflight` — what rules apply?
3. Classify: Is this a **code bug**, **environment drift**, or **snapshot parity** issue?
   - Code bug → Phase 2
   - Environment drift → `ci-parity` skill
   - Snapshot parity → `ci-parity` snapshot sync recipe

## Phase 2: Reproduce (TDD)
1. Write a failing test that captures the symptom (link to `go-tdd` RED phase)
2. If E2E: `go run ./cmd/verify browser` — capture the failure output
3. If CI-only: reproduce in Docker via `go run ./cmd/verify ci --with-docker`
4. Commit the reproduction: `test(debug): reproduce <symptom>`

## Phase 3: Fix → GREEN → Push
1. Apply minimal fix
2. Run the reproduction test — MUST pass
3. `go run ./cmd/verify precommit`
4. If UI: `go run ./cmd/verify browser`
5. Push via `./scripts/pushw.sh` and monitor with `gh run watch`

## Anti-patterns
- Do NOT hypothesize before reproducing. Write the test first.
- Do NOT run filtered Playwright subsets. Full suite or nothing.
- Do NOT skip Docker reproduction for CI-only failures.
