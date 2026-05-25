---
name: Go TDD Workflow
description: Execute the RED-GREEN-REFACTOR cycle for Go projects
triggers:
  - "writing tests"
  - "fixing bugs"
  - "implementing features"
  - "TDD"
  - "red green refactor"
mutating: true
---
# Go TDD Skill
## Session Start
> **Loop-Budget & Fail-Fast Guard**: If you are using a reasoning model (Pro/Opus) and executing coding/TDD, you are capped at **max 3 serial execution loops (compiles/tests)**. If compilation or a unit test fails **2 times consecutively**, you MUST HALT immediately, output a structured clinical diagnosis, and await human direction.
1. Run `go run ./cmd/verify preflight` — review active rules
2. Run `go run ./cmd/verify check-gates` — check active gates

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
## REFACTOR Phase (Clean Up)
1. Run `go run ./cmd/verify critique --baseline=HEAD~1`.
   - You MUST compare the output to ensure the total number of issues (especially Duplication/Clone Groups) is LESS THAN OR EQUAL TO the baseline. If violations increased, you MUST revert or fix the regression before committing.
2. Run `go run ./cmd/verify heal` — auto-fix structural issues
3. Run `go test ./path/to/package/` — confirm nothing broke
4. Stage and commit: `git add . && git commit -m "refactor(scope): clean up X"`
## Anti-Patterns (from Strict Lessons)
- Do NOT use `t.Parallel()` with `os.Setenv/os.Unsetenv` — causes flaky CI
- Do NOT lower coverage thresholds in `.agents/coverage.json`
- Do NOT skip `check-gates` — it enforces RED-before-GREEN ordering
