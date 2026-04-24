---
description: Process failure and feedback to codify corrections directly into repository rules.
---

# /learn Workflow

This workflow is triggered when a mistake is identified and a correction is required to prevent its recurrence. All learning MUST be codified into existing project artifacts to avoid unnecessary paperwork.

`[/learn: <mistake and required correction>]`

## Process

When the user issues a `/learn` command, you MUST execute the following steps:

1. **Analyze the Correction**:
   - Determine if the mistake is related to **Process, Style or TDD** (how code is written/tested).
   - Determine if the mistake is an **Architectural or Design** error (boundary violations, service coupling, etc.).
   - **De-duplication Check**: Scan `coding-standards.md`, existing skills in `.agents/skills/`, and existing verify subcommands for overlap.
   - **Determinism Test**: Can this correction be verified by running a command with a deterministic pass/fail outcome (no human judgment required)?
     - **If YES**: Skip step 2. Go directly to step 2.7 (Create Tool).
   - **Procedure Test**: Does this correction involve 2+ sequential steps that must be executed in order?
     - **If YES**: This is a Skill, not a Lesson. Go to step 2 option (c).

2. **Codify the Correction** (choose ONE):

   a. **For Declarative Rules (single constraint, no sequence)** — Process/Style/TDD:
      - Append under the appropriate ### subsection in coding-standards.md. Include [TRIGGER:].

   b. **For Architecture/Design** — if it changes a core principle:
      - Create ADR in `docs/adr/YYYY-MM-DD-[lesson].md`.
      - Link in `AGENTS.md` if it changes a global constraint.

   c. **For Procedural Patterns (multi-step sequence, checklist)** — Skill:
      - Check if an existing skill in `.agents/skills/` covers this domain.
      - If yes: update the existing SKILL.md with the new steps or failure pattern.
      - If no: create `.agents/skills/<name>/SKILL.md` with YAML frontmatter.
      - Add the skill to the table in `AGENTS.md` under `## SKILLS`.
      - Register in `.agents/verify-manifest.yaml` under `skills:`.

2.7. **Create Tool (Deterministic Check)**:
   - Create a `verify` subcommand in `cmd/verify/misc.go`.
   - Implement the check in `internal/maintenance/<name>.go` with a test.
   - Register in `cmd/verify/main.go` and `.agents/verify-manifest.yaml`.
   - If this replaces an existing Strict Lesson, retire the lesson per step 2.5.


2.5. **Retirement Check (Lesson Lifecycle)**:
   - Before adding a new Strict Lesson, check if any EXISTING lessons in coding-standards.md can be retired.
   - A lesson is eligible for retirement if:
     a. A `verify` subcommand now enforces the same check deterministically (tool replaces prose)
     b. A test in the test suite explicitly covers the failure case the lesson describes
   - If a lesson is eligible: remove it from coding-standards.md and add a comment in the test or verify command referencing the retired lesson (e.g., `// Replaces Strict Lesson: HTTPS Awareness`).
   - **Goal**: Keep Strict Lessons count ≤ 25. If adding a new lesson would exceed 25, you MUST retire at least one.

3. **Verify and Commit**:
   - Ensure the updated rule or new ADR is correctly formatted.
   - Execute an atomic commit with the message: `chore(learn): codify correction for [short description of mistake]`.

4. **Dynamic Reload**:
   - Confirm to the user that the rule has been codified and committed.
   - Acknowledge that this rule is now active and will be loaded in future `[/build-feature]` sessions.