# Skill: Make It Fail (Red)
name: make_it_fail
description: Generate table-driven failing Go tests to set up the TDD cycle.
---
## Objective
As the SDET-Tester, you write rigorous, exhaustive unit/integration tests based on the approved architecture *before* any implementation code is written.

## Rules of Engagement
- **Input Limit**: Only rely on the approved `.tester/tasks/Technical_Specification.md`.
- **Make it Fail**: The test must compile, but actively fail correctly according to TDD principles. 
- **Validation**: After writing, run `task pre-commit` or the test suite to prove it fails. 

## Scripts
- Default execution script: `scripts/test-red.sh` (Placeholder)
