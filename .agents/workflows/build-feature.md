---
description: Start the Autonomous AI Developer Pipeline sequence with a new feature idea
---

When the user types `/build-feature <idea>`, orchestrate the development process strictly using `.agents/config.yaml` personas and `.agents/personas/`.

## Pipeline State Machine

```
@ProductOwner → @SystemsArchitect → [Approved] → @SDET-Tester
@SDET-Tester → [RED Pass] → @BackendEngineer
@BackendEngineer → [GREEN Pass] → @ChiefCritic (⚠️ Pro model required)
@ChiefCritic → [P0/P1/P2 Defects] → @SDET-Tester
@ChiefCritic → [Clean / P3 only] → @BackendEngineer (REFACTOR)
@BackendEngineer → [REFACTOR Pass] → @SecurityEngineer
@SecurityEngineer → [P0/P1/P2 Defects] → @SDET-Tester
@SecurityEngineer → [Clean] → @ChaosMonkey
@ChaosMonkey → [Brittle] → @SDET-Tester
@ChaosMonkey → [Resilient] → ✅ Done
```

### Phase 1: Architecture & Planning (Human Intervention Required)

1. Act as the **@ProductOwner** and execute the `design_architecture` skill using the `<idea>` to define the user value, spec, and "Why".
2. Shift context and act as the **@SystemsArchitect** to review the spec for technical feasibility and refine the architecture.
   *(Wait for the user to explicitly approve the spec. Loop until approved.)*
   - The `implementation_plan.md` MUST include a **File Change Manifest**: an ordered table of every file to be created or modified, labeled `[NEW]`, `[MODIFY]`, or `[DELETE]`, with an estimated **token count** per file. Any file with >500 estimated tokens must be flagged **`[EXPENSIVE]`** and annotated with named sub-sections (e.g., "structs → handler logic → helpers"). Order the manifest smallest-first.
3. **@ProductOwner** authors three output files before autonomous execution begins:
   - `implementation_plan.md` — technical spec for Phase 2 personas. 
   - `.tester/tasks/vibe_check.md` — manual UX/aesthetic checklist generated from `.agents/vibe-check-template.md`. User/PO must check off all items.
   - `.tester/tasks/CHAOS_BRIEF.md` — Phase 3 chaos brief for **@ChaosMonkey**. Fill in the sabotage targets specific to this feature using the template at `.tester/tasks/CHAOS_BRIEF.md`.

### Phase 2: Autonomous Execution Loop (Multi-Conversation)

> [!IMPORTANT]
> Each step MUST run in a **NEW conversation window** (New Window Protocol) to minimize context cost and hallucination risk.
> After completing your turn, generate a handoff and instruct the user to open a new chat.

1. **@SDET-Tester**: Execute the `make_it_fail` skill to generate RED test cases.
   - *Transition*: `./scripts/agent-exec.sh set-phase RED`
   - *Handoff*: `./scripts/agent-exec.sh handoff BackendEngineer`
2. **@BackendEngineer**: Run `/resume`, then execute `make_it_pass` to pass the tests.
   - *Transition*: `./scripts/agent-exec.sh set-phase GREEN`
   - *Handoff*: `./scripts/agent-exec.sh handoff ChiefCritic`
3. **@ChiefCritic** *(⚠️ requires Pro model)*: Run `/resume`, then execute `critique_product`. Apply severity threshold from `defect_policy`. Only P0/P1/P2 defects trigger kick-back.
   - *Handoff*: `./scripts/agent-exec.sh handoff BackendEngineer` (if P0/P1/P2 found) or `handoff BackendEngineer` (REFACTOR if clean)
4. **@BackendEngineer**: Run `/resume`, then execute `make_it_better`.
   - *Transition*: `./scripts/agent-exec.sh set-phase REFACTOR`
   - *Handoff*: `./scripts/agent-exec.sh handoff SecurityEngineer`
5. **@SecurityEngineer**: Run `/resume`, then execute `audit_security`. Apply severity threshold from `defect_policy`. Only P0/P1/P2 findings trigger kick-back.
   - *Finalization*: `./scripts/agent-exec.sh handoff ChaosMonkey`

### Phase 3: Resilience & Chaos (Human Intervention Required)

6. **@ChaosMonkey**: Open a NEW chat. Run `/resume`. Read `.tester/tasks/CHAOS_BRIEF.md` (your mission brief). Execute sabotage targets in order. Report against success/failure conditions in the brief. **Never write fixes** — kick back to @SDET-Tester on failure.
7. **Final Verification**: @SystemsArchitect confirms all gates are green: `./scripts/agent-exec.sh status`. Run `./scripts/agent-exec.sh verify vibe-check` to ensure the human PO has validated the UX checklist. Run `/janitor` if workspace entropy is high.