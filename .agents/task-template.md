# Task: [TASK_NAME]

## Context & Objective
[Provide a brief description of the goal and any background information.]

## Persona & Phase
- **Active Persona**: [SystemsArchitect | SDET-Tester | BackendEngineer | etc.]
- **Current Phase**: [PLANNING | EXECUTION | VERIFICATION]
- **State Machine Phase**: [IDLE | RED | GREEN | REFACTOR]

## State Machine Gates
> **Rule**: All gates must pass for the feature to be considered complete.

| Gate ID | Description | Status | Verification Tool |
| :--- | :--- | :--- | :--- |
| **G.1: Red Test** | Prove feature absence (Fail) | [ ] | `scripts/agent-gate.sh red-test` |
| **G.2: API Spec** | Defined in `docs/api.md` | [ ] | Manual Review |
| **G.3: Implementation** | Unit tests pass (Green) | [ ] | `scripts/agent-gate.sh implementation` |
| **G.4: Lint/Coverage** | No regressions, ≥ threshold in `.agent/coverage-threshold` | [ ] | `scripts/pre-commit.sh` |
| **G.5: Browser** | UI/UX Native verification | [ ] | `browser_subagent` |

## Technical Brand Enforcement
- [ ] Aesthetic Check: [Sharp Editorial Earth]
- [ ] Typography: [Serif Headers / Uppercase Micro-copy]
- [ ] Surface: [Glassmorphism / Zero-radius]

## Loop Roadmap
| Step | Activity | Phase | Status |
| :--- | :--- | :--- | :--- |
| 1 | Research & Planning | PLANNING | [ ] |
| 2 | Red Test Implementation | EXECUTION | [ ] |
| 3 | Core Logic (Green) | EXECUTION | [ ] |
| 4 | Programmatic Verification | VERIFICATION | [ ] |
| 5 | Browser Native Validation | VERIFICATION | [ ] |

## Browser Verification Recordings
> **Mandatory**: Embed recordings of successful UI flows here.

- **Recording 1**: `verify_[feature]_ui_flow` - [Link](file:///Users/johnnyblase/gym/agbalumo/@logs/recordings/...)
- **Recording 2**: `verify_[feature]_edge_cases` - [Link](file:///Users/johnnyblase/gym/agbalumo/@logs/recordings/...)

## Artifact Logs
- [ ] `implementation_plan.md`
- [ ] `walkthrough.md`
- [ ] `task.md`
