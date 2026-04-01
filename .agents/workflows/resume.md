---
description: Bootstrap a new persona from a HANDOFF.md file to minimize context cost
---

# /resume Workflow

Use this workflow to initialize a new conversation window with the correct project state and persona.

1. **Read Project Baseline**:
   - `view_file .tester/tasks/HANDOFF.md` <!-- Identify the Target Persona and current phase -->
   - `view_file .tester/tasks/progress.md` <!-- Historical context of what's done -->
   - `view_file implementation_plan.md` <!-- The technical spec for the next steps -->

1.5. **Confirm Identity** *(tripwire — do this before any other action)*:
   - Read the `target_persona:` field from `HANDOFF.md`.
   - State your identity aloud in your **first response** to the user:
     > "I am **[target_persona]** as declared in HANDOFF.md. Feature: [feature]. Phase: [phase]."
   - **STOP condition**: If you cannot confirm your identity from `HANDOFF.md` (file missing, field absent, or feature name does not match your context), output the following and do nothing else:
     > "⛔ RESUME HALTED: Cannot confirm persona identity. Please verify HANDOFF.md exists and contains `target_persona`, `feature`, and `phase` fields, then retry /resume."

2. **Assume Persona**:
   - Update your internal state to match the declared persona role (e.g., if `target_persona: BackendEngineer`, switch to TDD Green phase mindset).
   - Load persona-specific instructions from `.agents/personas/<target_persona>.yaml`.

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

> [!IMPORTANT]
> Step 1.5 is the anti-hallucination tripwire. It must execute before any code, tool calls, or environment checks. A persona that skips it is operating without authorization.
