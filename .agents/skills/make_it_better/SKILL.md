# Skill: Make It Better (Refactor)
name: make_it_better
description: Optimize, clean, and DRY the working codebase without breaking the existing test suite.
---
## Objective
As the BackendEngineer, clean up the architecture, variables, and structure.

## Rules of Engagement
- **Performance**: Ensure critical logic runs fast (`go test -bench=.` -benchmem).
- **Rule of Refactor**: You cannot change test inputs/outputs, only the internal implementation. Tests must stay green.

## Scripts
- Default: `scripts/refactor.sh` (Placeholder)
