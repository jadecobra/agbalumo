# Implementation Plan - Context Cost Optimization

## 1. Problem Statement
The current codebase has an RMS (Root Mean Square) Context Cost of **192.26**. 
The top 5 most expensive files account for a significant portion of this cost, especially `critique_report.md` which is an outlier at 2651 lines. 
High LOC files increase cognitive load for the AI squad, lead to context-exhaustion, and increase the risk of hallucinations.

## 2. Top 5 Expensive Files & Proposed Remediation

| File | LOC | Type | Strategy |
| :--- | :--- | :--- | :--- |
| `critique_report.md` | 2651 | Report | **Exclude/Archive**: Move to `.agents/reports/` and update `cost.go` to ignore it. |
| `internal/agent/verify_apispec_test.go` | 641 | Test | **Modularize**: Split into `verify_apispec_routes_test.go` and `verify_apispec_harness_test.go`. |
| `internal/agent/security.go` | 627 | Logic | **Refactor**: Split into `security_web.go`, `security_sql.go`, and `security_fs.go`. |
| `internal/repository/sqlite/sqlite_listing_test.go` | 578 | Test | **Split**: Extract CRUD and Search tests into separate files. |
| `internal/agent/verify_test.go` | 569 | Test | **Split**: Break down by gate type (RedTest, Coverage, etc). |

## 3. Technical Requirements
- **RMS Target**: Reduce RMS LOC below **110**.
- **Constraint**: No reduction in test coverage.
- **Constraint**: All verification gates must pass.
- **Protocol**: Standard Red-Green-Refactor logic for splitting logic.

## 4. Phase 1: Architecture (Planning)
- **@ProductOwner**: Establish the "Why" (done).
- **@SystemsArchitect**: Refine the split points for `security.go` and `verify_apispec_test.go`.
- **Approved**: Human confirmation required.

## 5. Phase 2: Execution (Autonomous)
- **@SDET-Tester**: Identify critical paths in `security.go` that must not break during split.
- **@BackendEngineer**: Perform the splits.
- **@ChiefCritic**: Verify the new structure didn't introduce technical debt.

## 6. Chaos Brief Targets
- **ChaosMonkey**: Sabotage the `cost` calculation by adding hidden long files.
- **ChaosMonkey**: Introduce a circular dependency during `security.go` split.
