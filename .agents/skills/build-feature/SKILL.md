---
name: "Build Feature"
description: "Execute the end-to-end engineering lifecycle for a new feature across planning, TDD execution, resilience audit, and deployment."
triggers:
  - "/build-feature"
  - "build feature"
  - "implement feature"
  - "new feature"
mutating: true
---

# Build Feature Skill

Execute the entire engineering lifecycle for a new feature in a continuous workflow.
Phase 1 requires explicit user approval of the task plan.
Phase 2 and Phase 3 run autonomously once approved.
Checkpoints are recorded through git commits.

## Phase 1: Product Architecture and Planning
1. Review `.agents/skills/flash-plan/SKILL.md` for session start rules and the Architect Protocol.
2. Initialize `task.md` containing a Decision Log section and an Execution Plan section.
3. Record every architectural decision, deleted abstraction, and latency tradeoff in the Decision Log.
4. Populate the Execution Plan with markdown checklist items for Phase 2 and Phase 3.
5. Stop and present the plan to the user. Do not proceed to Phase 2 until the user explicitly approves `task.md`.

## Phase 2: Autonomous Execution Loop (TDD)
1. Verify model budget guardrails. If using an expensive reasoning model, verify the user provided the OVERRIDE flag before starting iterative execution loops.
2. Follow the RED-GREEN-REFACTOR protocol in `.agents/skills/go-tdd/SKILL.md`.
3. Consult the testing section in `.agents/coding-standards.md`.
4. Run `go run ./cmd/verify precommit` before committing changes.
5. If HTML, CSS, or UI templates changed, run `go run ./cmd/verify browser` locally to detect visual regressions.

## Phase 3: Audit and Resilience
1. Run security tooling via `gosec ./...`. Resolve any high or medium defects immediately.
2. Verify SQLite query execution plans using `EXPLAIN QUERY PLAN`.
3. Execute tests against a file-backed SQLite database to verify write-ahead logging and concurrency behavior.
4. Verify external and template contracts using `go run ./cmd/verify template-drift` and `go run ./cmd/verify api-spec`.
5. For UI changes, capture a screenshot of the completed flow using the browser subagent.
6. Push changes using `./scripts/pushw.sh` and monitor results with `gh run watch`.
7. Review `task.md` for repeated multi-step procedures. Extract repeated patterns into new skills or deterministic verify commands.
