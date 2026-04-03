---
description: Bootstrap a new persona from a HANDOFF.md file to minimize context cost
---

# Resume Workflow

Use this workflow when starting a **NEW conversation window** (New Window Protocol) after a persona handoff. This bootstraps your memory and current task state without requiring the full conversation history.

## 0. Initial Bootstrap

When you see the user type `/resume`, immediately parse the latest **HANDOFF.md** and **progress.md** files.

// turbo
1. Run `ls -t .tester/tasks/HANDOFF.md progress.md | head -n 2` to confirm latest state files.
// turbo
2. View the contents of the latest **HANDOFF.md**.

## 1. Persona Alignment

Identity your current persona as defined in the `HANDOFF.md`:
- `@SDET-Tester`: Owns the RED phase (Proof).
- `@BackendEngineer`: Owns the GREEN and REFACTOR phases (Fix).
- `@SecurityEngineer` / `@ChiefCritic`: Owns the Audit phase (Verification).

## 2. Environment Setup

// turbo
1. Run `./scripts/agent-exec.sh status` to verify current feature and phase.
// turbo
2. Verify you have the correct file context for the current task.

## 3. Execution

Continue the task defined in the `logic_brief` section of the `HANDOFF.md`.
- **RED**: Write failing tests. Verify with `./scripts/agent-exec.sh verify red-test`.
- **GREEN**: Pass tests. Verify with `./scripts/agent-exec.sh verify implementation`.
- **REFACTOR**: Optimize. Verify with `task lint` and `./scripts/agent-exec.sh verify coverage`.

## 4. Final Verification and Handoff

Before completing your session, ensure you follow the handoff protocol:
1. Pass all required gates for your phase.
2. Update the `progress.md` via `pending_update.md`.
3. Run `./scripts/agent-exec.sh handoff <next_persona>` to generate the NEW `HANDOFF.md`.
4. Instruct the user to open a NEW chat window.
