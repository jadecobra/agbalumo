---
description: Bootstrap a new persona from a HANDOFF.md file to minimize context cost
---

# /resume Workflow

Use this workflow to initialize a new conversation window with the correct project state and persona.

1. **Read Project Baseline**:
   - `view_file .tester/tasks/HANDOFF.md` <!-- Identify the Target Persona and current phase -->
   - `view_file .tester/tasks/progress.md` <!-- Historical context of what's done -->
   - `view_file implementation_plan.md` <!-- The technical spec for the next steps -->

2. **Assume Persona**:
   - Confirm your identity based on the `Target Persona` field in `HANDOFF.md`.
   - Update your internal state (e.g., if you are now the `BackendEngineer`, switch to the TDD Green phase).

3. **Validate Environment**:
   - Run `harness status --text` to confirm the CLI state matches the handoff file.
   - Run `task lint` and `task test` to ensure the workspace is clean.

4. **Initialize Execution**:
   - Summarize your understanding of the next task to the user.
   - Proceed with the next step of the pipeline (e.g. implementing the logic to pass the RED tests).

5. **Cleanup**:
   - Once initialized, you may delete `.tester/tasks/HANDOFF.md` to prevent accidental re-resumes.
     - `run_command rm .tester/tasks/HANDOFF.md`

> [!TIP]
> This command is designed to be run as the **FIRST** interaction in a brand-new chat window.
