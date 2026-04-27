---
description: Execute the end-to-end Engineering pipeline for a new feature.
---

# Engineering Workflows

`/build-feature <idea>`
When the user types `/build-feature <idea>`, act as a Senior Product Engineer. Execute the entire lifecycle (Architecture, TDD, Security, Resilience, and Observability) in a single continuous session.

**Phase 1 is complete ONLY when the user approves the updated `task.md`. Once approved, jump to Phase 2 and do not ask for permission again. Your Git commits are your checkpoints**.

## Phases
| Phase | Description | File |
|---|---|---|
| Phase 1 | Product Architecture & Planning | [build-feature-phase1.md](file:///Users/johnnyblase/gym/agbalumo/.agents/workflows/build-feature-phase1.md) |
| Phase 2 | Autonomous Execution Loop (TDD) | [build-feature-phase2.md](file:///Users/johnnyblase/gym/agbalumo/.agents/workflows/build-feature-phase2.md) |
| Phase 3 | Audit & Resilience | [build-feature-phase3.md](file:///Users/johnnyblase/gym/agbalumo/.agents/workflows/build-feature-phase3.md) |