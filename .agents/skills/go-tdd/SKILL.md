---
name: Go TDD Workflow
description: Execute the RED-GREEN-REFACTOR cycle for Go projects
triggers:
  - "writing tests"
  - "fixing bugs"
  - "implementing features"
  - "TDD"
  - "red green refactor"
mutating: false
---
# Go TDD Skill
## Session Start
1. Run `go run ./cmd/verify preflight` — review active rules
2. Run `go run ./cmd/verify check-gates` — detect current TDD phase
## RED Phase (Write Failing Test)
1. Create or modify `*_test.go` file with the new test case
2. Use table-driven tests (see `.agents/workflows/coding-standards.md` → Testing Conventions)
3. Check `internal/testutil/` for existing helpers before writing custom setup
4. Run: `go test -run TestNewFeature ./path/to/package/`
5. **MUST** see test FAIL (exit code 1)
6. Stage and commit: `git add *_test.go && git commit -m "test(scope): add failing test for X"`
## GREEN Phase (Make Test Pass)
1. Write minimum implementation to pass the test
2. Run: `go test -run TestNewFeature ./path/to/package/`
3. If FAIL after 3 attempts:
   - HALT
   - Read the raw traceback
   - Hypothesize why (wrong mock? missing interface method? SQL schema?)
   - Present hypothesis to user — WAIT for guidance
4. Stage and commit: `git add . && git commit -m "feat(scope): implement X"`
## REFACTOR Phase (Clean Up)
1. Run `go run ./cmd/verify critique` — check for violations
2. Run `go run ./cmd/verify heal` — auto-fix structural issues
3. Run `go test ./path/to/package/` — confirm nothing broke
4. Stage and commit: `git add . && git commit -m "refactor(scope): clean up X"`
## Anti-Patterns (from Strict Lessons)
- Do NOT use `t.Parallel()` with `os.Setenv/os.Unsetenv` — causes flaky CI
- Do NOT lower coverage thresholds in `.agents/coverage.json`
- Do NOT skip `check-gates` — it enforces RED-before-GREEN ordering
