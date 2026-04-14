---
description: Process failure and feedback to codify corrections directly into repository rules.
---

# /learn Workflow

This workflow is triggered when a mistake is identified and a correction is required to prevent its recurrence. All learning MUST be codified into existing project artifacts to avoid unnecessary paperwork.

`[/learn: <mistake and required correction>]`

## Process

When the user issues a `/learn` command, you MUST execute the following steps:

1. **Analyze the Correction**:
   - Determine if the mistake is related to **Process, Style, or TDD** (how code is written/tested).
   - Determine if the mistake is an **Architectural or Design** error (boundary violations, service coupling, etc.).

2. **Codify the Correction**:

   - **For Process/Style/TDD**:
     - Append the corrected rule directly to the bottom of [.agents/workflows/coding-standards.md](file:///Users/johnnyblase/gym/agbalumo/.agents/workflows/coding-standards.md) under the `# Strict Lessons` section.
     - Use a clear, imperative bullet point (e.g., "* The agent MUST always...").

   - **For Architecture/Design**:
     - Create a formal Architecture Decision Record (ADR) in `docs/adr/YYYY-MM-DD-[lesson].md`.
     - Use the template at [docs/adr/template.md](file:///Users/johnnyblase/gym/agbalumo/docs/adr/template.md).
     - Link the new ADR in `AGENTS.md` if it changes a global architectural constraint.

3. **Verify and Commit**:
   - Ensure the updated rule or new ADR is correctly formatted.
   - Execute an atomic commit with the message: `chore(learn): codify correction for [short description of mistake]`.

4. **Dynamic Reload**:
   - Confirm to the user that the rule has been codified and committed.
   - Acknowledge that this rule is now active and will be loaded in future `[/build-feature]` sessions.
