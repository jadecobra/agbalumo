---
description: Workflow for implementing new features with strict TDD
---

# Feature Implementation Workflow

This workflow enforces the "Test First" rule to prevent breaking changes.

1.  **Plan**: Define the feature and expected behavior.
    -   *Action*: Update `task.md` or Create `implementation_plan.md`.
    -   *Check*: Does this require a schema change? If so, verify migration plan.

2.  **Test (Red)**: Write a failing test case.
    -   *Role*: SDET Agent
    -   *Command*: `go test ./internal/package_name` (Should FAIL)
    -   *Verify*: Ensure the failure is due to *missing logic*, not compilation error (unless adding new API).

3.  **Implement (Green)**: Write the minimal code to pass the test.
    -   *Role*: Backend Agent
    -   *Constraint*: Do not add features not covered by the test.
    -   *Command*: `go test ./internal/package_name` (Should PASS)

4.  **Refactor**: Clean up the code while keeping tests passing.
    -   *Role*: Backend Agent / Lead Architect
    -   *Command*: `go test ./...` (Regression check)

5.  **Secure**: Verify security implications.
    -   *Role*: Security Engineer
    -   *Check*: Input validation, headers, permissions.

6.  **Verify**: Final full-suite check.
    -   *Command*: `./scripts/pre-commit.sh`
