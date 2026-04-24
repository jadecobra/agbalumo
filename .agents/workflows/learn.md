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
   - **De-duplication Check**: Scan `coding-standards.md` for existing rules related to this mistake.

2. **Codify the Correction**:
   - **Merge**: If a related rule exists, refactor it into a single, more robust abstraction.
   - **For Process/Style/TDD**:
     - Append the corrected rule under the appropriate ### subsection (CI & Infrastructure, UI & Frontend, Security & Environment, or Testing) within the # Strict Lessons section of [.agents/workflows/coding-standards.md](file:///Users/johnnyblase/gym/agbalumo/.agents/workflows/coding-standards.md). Include a [TRIGGER: ...] annotation.
     - Use a clear, imperative bullet point (e.g., "* The agent MUST always...").

   - **For Architecture/Design**: if it changes a core principle
     - Create a formal Architecture Decision Record (ADR) in `docs/adr/YYYY-MM-DD-[lesson].md`.
     - Use the template at [docs/adr/template.md](file:///Users/johnnyblase/gym/agbalumo/docs/adr/template.md).
     - Link the new ADR in `AGENTS.md` if it changes a global architectural constraint.

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