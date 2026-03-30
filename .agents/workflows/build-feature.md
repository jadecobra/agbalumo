---
description: Start the Autonomous AI Developer Pipeline sequence with a new feature idea
---
When the user types `/build-feature <idea>`, orchestrate the development process strictly using `.agents/config.yaml` personas and `.agents/personas/`.

### Phase 1: Architecture & Planning (Human Intervention Required)

1. Act as the **ProductOwner** and execute the `design_architecture` skill using the `<idea>` to define the user value, spec, and "Why".
2. Shift context and act as the **SystemsArchitect** to review the spec for technical feasibility and refine the architecture. 
   *(Wait for the user to explicitly approve the spec. If the user provides feedback or adds comments directly to the Markdown file, act as the ProductOwner/SystemsArchitect again to re-read and revise the document. Loop this step until they approve).*

### Phase 2: Autonomous Execution Loop (Fully Automated)

> [!NOTE]
> Steps 3 through 7 are now fully automated by the Go harness via the `antigravity chat` CLI. Once Phase 1 is approved, the human-led agent triggers `./scripts/agent-exec.sh init` and the harness autonomously orchestrates the following personas.

3. **SDET-Tester**: Execute the `make_it_fail` skill to generate rigorous RED test cases.
4. **BackendEngineer**: Execute the `make_it_pass` skill to pass the tests (GREEN phase).
5. **ChiefCritic**: Execute the `critique_product` skill to audit the implementation. If flaws are found, the harness returns to Step 4.
6. **BackendEngineer**: Execute the `make_it_better` skill to refactor the working code optimally.
7. **SecurityEngineer**: Execute the `audit_security` skill to review for OWASP vulnerabilities.

### Phase 3: Resilience & Chaos (Human Intervention Required)

8. **Chaos Stress Test**: Shift context to the **ChaosMonkey** and execute a "Hard Mode" challenge. Intentionally inject a state failure or test sabotage. If the **SystemsArchitect** cannot restore the environment or the **SDET-Tester** fails to detect the sabotage, return to Step 6.
9. Final Verification and Squad Sync.
