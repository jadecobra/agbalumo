---
description: Start the Autonomous AI Developer Pipeline sequence with a new feature idea
---
When the user types `/build-feature <idea>`, orchestrate the development process strictly using `.agents/config.yaml` personas and `.agents/personas/`.

### Phase 1: Architecture & Planning (Human Intervention Required)

1. Act as the **ProductOwner** and execute the `design_architecture` skill using the `<idea>` to define the user value, spec, and "Why".
2. Shift context and act as the **SystemsArchitect** to review the spec for technical feasibility and refine the architecture. 
   *(Wait for the user to explicitly approve the spec. If the user provides feedback or adds comments directly to the Markdown file, act as the ProductOwner/SystemsArchitect again to re-read and revise the document. Loop this step until they approve).*

### Phase 2: Autonomous Execution Loop (Multi-Conversation)

> [!IMPORTANT]
> To minimize context cost and hallucination risk, each step below MUST be performed in a **NEW conversation**. 
> After completing your turn, generate a handoff for the next persona and instruct the user to start a new chat.

1. **SDET-Tester**: Execute the `make_it_fail` skill to generate RED test cases.
   - *Handoff*: `harness handoff BackendEngineer`
2. **BackendEngineer**: Run `/resume`, then execute `make_it_pass` to pass the tests.
   - *Handoff*: `harness handoff ChiefCritic`
3. **ChiefCritic**: Run `/resume`, then execute `critique_product`. If flaws are found, handoff back to `BackendEngineer`.
   - *Handoff*: `harness handoff BackendEngineer` or `harness handoff SecurityEngineer`
4. **BackendEngineer**: Run `/resume`, then execute `make_it_better`.
   - *Handoff*: `harness handoff SecurityEngineer`
5. **SecurityEngineer**: Run `/resume`, then execute `audit_security`.
   - *Finalization*: `harness handoff SystemsArchitect`

### Phase 3: Resilience & Chaos (Human Intervention Required)

8. **Chaos Stress Test**: Shift context to the **ChaosMonkey** and execute a "Hard Mode" challenge. Intentionally inject a state failure or test sabotage. If the **SystemsArchitect** cannot restore the environment or the **SDET-Tester** fails to detect the sabotage, return to Step 6.
9. Final Verification and Squad Sync.
