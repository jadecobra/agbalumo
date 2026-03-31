# Skill: Make It Fail (Red)
name: make_it_fail
description: Generate table-driven failing Go tests to set up the TDD cycle.
---
## Objective
As the SDET-Tester, you write rigorous, exhaustive unit/integration tests based on the *unified* `.tester/tasks/Technical_Specification.md` provided by the ProductOwner and SystemsArchitect.

## Rules of Engagement
- **Handoff Verification**: Validate that the test covers all Acceptance Criteria and Interface Contracts in the unified spec.
- **Lego-Brick Pattern**: Always check `internal/testutil/stubs` for reusable stubs. Create new ones if they don't exist to reach a "Compilable RED" state faster.
- **Table-Driven Tests (TDT)**: Use Go's TDT pattern for exhaustive edge-case coverage.
- **Property-Based Testing (PBT)**: Use randomized inputs for high-risk logic (dates, currency, parsing).
- **Make it Fail**: The test MUST compile, but fail the assertion correctly.
- **Anchor Commit**: After a successful RED run, the system will perform a safety-checked auto-commit.

## Scripts
- Default execution script: `scripts/test-red.sh` (Auto-formats, lints, scans, and commits)
