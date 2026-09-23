---
name: Go TDD Workflow
description: Execute the RED-GREEN-REFACTOR cycle for Go projects, surgical hotfixes, refactoring scans, and test-driven debugging.
triggers:
  - "writing tests"
  - "fixing bugs"
  - "implementing features"
  - "TDD"
  - "red green refactor"
  - "/hotfix"
  - "hotfix"
  - "/refactor"
  - "refactor"
  - "/debug"
  - "debug"
mutating: true
---
# Go TDD Skill
## Session Start
> **Loop-Budget & Fail-Fast Guard**: If you are using a reasoning model (Pro/Opus) and executing coding/TDD, you are capped at **max 3 serial execution loops (compiles/tests)**. If compilation or a unit test fails **2 times consecutively**, you MUST HALT immediately, output a structured clinical diagnosis, and await human direction.
1. Run `go run ./cmd/verify preflight` to review active rules.
2. Run `go run ./cmd/verify check-gates` to check active gates.

## RED & GREEN Phase (Local TDD Loop)
To minimize token consumption, reduce serial git overhead, and maintain TDD discipline:
1. **Write the Failing Test**: Create or modify `*_test.go` with table-driven tests (using `internal/testutil/` helpers where possible).
2. **Local RED Verification**: Run `go test -run TestNewFeature ./path/to/package/`. You **MUST** see the test fail (exit code 1) in your local terminal output. Capture this output in your context as proof of RED state.
3. **Write Implementation**: Write the minimal code to satisfy the failing test.
4. **Local GREEN Verification**: Run `go test -run TestNewFeature ./path/to/package/` again.
5. If GREEN fails after 2 attempts:
   - HALT immediately.
   - Do NOT continue making guessing edits.
   - Present a structured diagnosis of the traceback to the user with hypotheses.
6. **Atomic Commit**: Once the test passes locally (GREEN), stage and commit both the test and implementation files together as a single atomic commit: `feat(scope): implement X with unit tests` or `fix(scope): resolve Y with unit tests`.

## Surgical Hotfix Path (/hotfix)
When triggered via `/hotfix <description>`:
1. Confirm failing test exists or write reproduction test first.
2. Write minimal fix to make the test pass.
3. Run `go run ./cmd/verify precommit`. If UI was modified, run `go run ./cmd/verify browser`.
4. Push via `./scripts/pushw.sh` and monitor `gh run watch --exit-status`.
5. Skip Phase 1 planning, Phase 3 chaos tests, and ADR generation. Commit: `fix(<scope>): <description>`.

## Debug & Triage Path (/debug)
When triggered via `/debug <symptom>`:
1. Scope recent changes via `git diff --name-only HEAD~3 HEAD`.
2. Classify failure domain:
   - Code bug: write reproduction test (RED) per Go TDD loop.
   - Environment drift: follow `ci-parity` skill.
   - Snapshot parity: run `go run ./cmd/verify snapshot-parity` and snapshot sync recipe.
3. For CI-only failures: reproduce locally in Docker via `go run ./cmd/verify ci --with-docker`.

## REFACTOR Phase (Clean Up & /refactor)
1. Run `go run ./cmd/verify critique --baseline=HEAD~1`.
   - You MUST compare the output to ensure the total number of issues (especially Duplication/Clone Groups) is LESS THAN OR EQUAL TO the baseline. If violations increased, you MUST revert or fix the regression before committing.
2. Run `go run ./cmd/verify heal` to auto-fix structural issues.
3. Run `go test ./path/to/package/` to confirm nothing broke.
4. Stage and commit: `git add . && git commit -m "refactor(scope): clean up X"`.
5. **Proactive Improvement Scan**: On modified packages, evaluate:
   - `go run ./cmd/verify context-cost` (file size and token density).
   - `go run ./cmd/verify deprecated` (new migration opportunities).
   - `go run ./cmd/verify agents-coverage` (missing package AGENTS.md).
   - `go run ./cmd/verify critique --baseline HEAD~1` (complexity score change).
   Present at most 3 tool-grounded improvement suggestions referencing file and line numbers with verify tool output. Wait for user approval before implementing suggestions.

## Anti-Patterns (from Strict Lessons)
- Do NOT use `t.Parallel()` with `os.Setenv/os.Unsetenv` because it causes flaky CI.
- Do NOT lower coverage thresholds in `.agents/coverage.json`.
- Do NOT skip `check-gates` because it enforces RED-before-GREEN ordering.
- Do NOT replace implicit SQLite rowid tiebreakers (e.g. `rowid ASC`) with explicit string ID tiebreakers (e.g. `id ASC`) in pagination or default ordering clauses unless explicitly required by domain specifications. Secondary indexes on rowid tables implicitly append `rowid ASC`; substituting string IDs alters natural insertion-order sorting and breaks index-order dependent E2E assertions.
