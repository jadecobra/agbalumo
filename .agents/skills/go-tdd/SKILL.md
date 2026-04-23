---
name: Go TDD Workflow
description: Execute the RED-GREEN-REFACTOR cycle for Go projects
---

# Go TDD Skill

## Pre-conditions

1. Run `go run ./cmd/verify preflight` to load active rules
2. Run `go run ./cmd/verify check-gates` to detect current phase

## RED Phase

1. Write test file with `_test.go` suffix
2. Run `go test -run TestNewFeature ./path/to/package`
3. Verify test FAILS (exit code 1)
4. `git add *_test.go && git commit -m "test(scope): add failing test for X"`

## GREEN Phase

1. Write minimum implementation
2. Run `go test -run TestNewFeature ./path/to/package`
3. If FAIL after 3 attempts: HALT, read traceback, hypothesize
4. `git add . && git commit -m "feat(scope): implement X"`

## REFACTOR Phase

1. Run `go run ./cmd/verify critique`
2. Fix violations
3. Run `go run ./cmd/verify heal`
4. `git commit -m "refactor(scope): clean up X"`
