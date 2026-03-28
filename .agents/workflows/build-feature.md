---
description: Start the Autonomous AI Developer Pipeline sequence with a new feature idea
---
When the user types `/build-feature <idea>`, orchestrate the development process strictly using `.agents/agent.yaml` personas and `.agents/skills/`.

### Execution Sequence:

1. Act as the **LeadArchitect** and execute the `design_architecture` skill using the `<idea>`.
   *(Wait for the user to explicitly approve the spec. If the user provides feedback or adds comments directly to the Markdown file, act as the LeadArchitect again to re-read and revise the document. Loop this step until they approve).*
2. Shift context, act as the **SDET-Tester**, and execute the `make_it_fail` skill to generate rigorous RED test cases.
3. Shift context, act as the **BackendEngineer**, and execute the `make_it_pass` skill to pass the tests (GREEN phase).
4. Shift context, act as the **ChiefCritic**, and execute the `critique_product` skill to tear down the implementation and find flaws. If flaws are found, return to step 3.
5. Shift context, act as the **BackendEngineer**, and execute the `make_it_better` skill to refactor the working code optimally.
6. Shift context, act as the **SecurityEngineer**, and execute the `audit_security` skill to review for OWASP vulnerabilities.
