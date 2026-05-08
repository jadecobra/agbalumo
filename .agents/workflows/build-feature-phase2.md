---
description: Phase 2: Autonomous Execution Loop (TDD)
---
# Phase 2: Autonomous Execution Loop (TDD)

**Context Detection**: If a Flash Planning prompt preceded this, skip to Actions.

### Pre-conditions
- Read `.agents/skills/go-tdd/SKILL.md`.
- Read `.agents/workflows/coding-standards.md` → Testing section.

### Actions

1.  **RED Phase (Failing Test)**
    - Write a failing test in the appropriate `*_test.go` file.
    - Verify the test fails by running `go test ./path/to/package/`.
    - Commit the failing test: `git add . && git commit -m "test(scope): add failing test for X" --no-verify`.

2.  **GREEN Phase (Implementation)**
    - Write the minimum code required to make the test pass.
    - Verify the test passes by running `go test ./path/to/package/`.
    - If you fail to achieve GREEN after 3 attempts, HALT and seek guidance.
    - Commit the implementation: `git add . && git commit -m "feat(scope): implement X"`.

3.  **REFACTOR Phase (Optimization)**
    - Run `go run ./cmd/verify heal` to auto-fix structural issues.
    - Run `go run ./cmd/verify critique` and address any regressions.
    - Commit the cleanup: `git add . && git commit -m "refactor(scope): clean up X"`.

### Verification
- Run `go run ./cmd/verify precommit` before commit to ensure all gates pass.
