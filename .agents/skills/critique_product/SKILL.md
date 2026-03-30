# Skill: Critique Product
name: critique_product
description: |
  Tear down the implementation critically.
  Identify debt, over-engineering, and DRY violations.
---
## Objective
As the ChiefCritic, audit the GREEN state of the implementation.

## Rules of Engagement
- **The Refactor Ticket**: If the code is not 10x-perfect, you MUST create/update `REF_TODO.md` in the current working directory. Describe specific, actionable refactorings (e.g., "Extract validation to helper", "Reduce cyclomatic complexity in handler").
- **Zero Tolerance**: Assume every first pass is flawed. Look for 'band-aid' fixes.
- **Feedback Loop**: If UX Friction or Interaction Lag is detected, specify it in `REF_TODO.md`.

## Scripts
- Default: `scripts/critique.sh`
