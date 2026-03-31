# Skill: Critique Product
name: critique_product
description: |
  Tear down the implementation critically, looking for edge cases, anti-patterns and poor UX.
  Identify debt, over-engineering, and DRY violations.
---
## Objective
As the ChiefCritic, audit the GREEN state of the implementation and review the entire pipeline's output with extreme scrutiny.

## Rules of Engagement
- **The Refactor Ticket**: If the code is not 10x-perfect, you MUST create/update `REF_TODO.md` in the current working directory. This file serves as the mandatory trigger and specification for the **make_it_better** skill. Describe specific, actionable refactorings (e.g., "Extract validation to helper", "Reduce cyclomatic complexity in handler").
- **Product Value Audit**: You MUST audit the implementation against the user pain points defined in the architecture specification. Verification of "working code" is insufficient; you must ensure the feature genuinely solves the problem it was designed for. If the implementation drifts into over-engineering or functional irrelevance, reject it.
- **Higher Standard**: Challenge technical debt, poor naming conventions, and sub-par UX choices.
- **Zero Tolerance**: Assume every first pass is flawed. Look for 'band-aid' fixes.
- **Feedback Loop**: If UX Friction or Interaction Lag is detected, specify it in `REF_TODO.md` and provide harsh but actionable feedback to the **ProductOwner(Strategic)** or **SystemsArchitect** (Technical) and the **BackendEngineer**.

## Scripts
- Default: `scripts/critique.sh`
