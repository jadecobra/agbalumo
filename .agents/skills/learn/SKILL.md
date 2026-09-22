---
name: "Learn"
description: "Process mistakes and feedback to codify permanent corrections into repository rules, ADRs, skills, or verify tools."
triggers:
  - "/learn"
  - "learn"
  - "codify lesson"
  - "record mistake"
mutating: true
---

# Learn Skill

Codify mistakes and user feedback directly into repository assets to prevent recurrence without creating unnecessary paperwork.

## Phase 1: Analyze the Correction
1. Classify the mistake into one of three domains: process and style, architectural boundaries, or deterministic checks.
2. Check for duplicate coverage across `.agents/coding-standards.md`, existing skills, and verify commands.
3. Apply the determinism test. If the check can be verified by a command with an objective pass or fail result, create a verify tool.
4. Apply the procedure test. If the correction requires two or more sequential steps, define or update a skill.
5. Apply the risk and blast radius test. Verify that the proposed lesson or rule does not conflict with existing invariants, bypass essential safety rails, or create unintended secondary failures.

## Phase 2: Codify the Correction
Choose exactly one destination based on classification.

### A. Declarative Rules (Single Constraint)
Append the constraint under the relevant section in `.agents/coding-standards.md`.
Include a bracketed trigger tag such as `[TRIGGER: handler_change]`.

### B. Architectural Principles
Draft an Architecture Decision Record in `docs/adr/YYYY-MM-DD-<title>.md`.
Reference the ADR in `AGENTS.md` if it alters a global constraint.

### C. Procedural Workflows
Update an existing skill in `.agents/skills/` or create a new skill directory.
Register new skills in `.agents/skills/RESOLVER.md` and `.agents/verify-manifest.yaml`.

### D. Deterministic Checks (Verify Tool)
Follow `.agents/skills/verify-authoring/SKILL.md` to author a verify subcommand.
Implement the logic in `internal/maintenance/` accompanied by unit tests.
Register the command in `cmd/verify/` and `.agents/verify-manifest.yaml`.
If the new tool enforces an existing strict lesson, retire that prose lesson.

## Phase 3: Retirement and Ceilings
Check if existing strict lessons in `.agents/coding-standards.md` can be retired.
A lesson qualifies for retirement when a verify command enforces it deterministically or unit tests cover the failure mode.
Keep the active strict lessons count at or below 20.

## Phase 4: Verification and Commit
Run `go run ./cmd/verify skill-conformance` and `go run ./cmd/verify check-resolvable`.
Run `go run ./cmd/verify lessons-conformance`.
Commit the changes using `chore(learn): codify correction for <short description>`.
Confirm to the user that the lesson is committed and active.
