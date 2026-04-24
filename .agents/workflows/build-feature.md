---
description: Execute the end-to-end Engineering pipeline for a new feature.
---

Engineering Workflows

`/build-feature <idea>`
When the user types `/build-feature <idea>`, act as a Senior Product Engineer. Execute the entire lifecycle (Architecture, TDD, Security, Resilience, and Observability) in a single continuous session.

**Phase 1 is complete ONLY when the user approves the updated `task.md`. Once approved, jump to Phase 2 and do not ask for permission again. Your Git commits are your checkpoints**.

## Phase 1: Product Architecture & Planning (The Iterative Algorithm)

HALT and execute this protocol. Initialize `task.md` immediately to externalize state and preserve context window.

0. **Session Bootstrap (Deterministic)**:
   - Run `go run ./cmd/verify preflight` to load rules relevant to this session's file changes.
   - Read `.agents/invariants.json` for hardcoded project constants (port, protocol, DB engine).
   - Read `.agents/verify-manifest.yaml` to understand which verify subcommands to run at each workflow stage.
   - **Rule**: Do NOT proceed to architecture planning until preflight output has been reviewed.


1. **Initialize Decision Log**:

- Create `task.md` with two sections: `# Decision Log` (append-only rationale) and `# Execution Plan` (the Phase 2/3 checkboxes).
- **Rule**: Log every decision, deleted component, and architectural tradeoff here before moving to the next turn.

2. **The Product Interrogation (The Kill-Switch Gate)**:

- **The Objective**: Every feature must contribute to finding food in < 60 seconds.
- **The Protocol**: Ask one piercing question at a time to challenge the core assumption.
- **Complexity Kill-Switch**: You MUST challenge the user to prove that the proposed complexity doesn't increase "Time to First Result." If it adds UI steps or DB latency without a 2x increase in utility, you MUST propose deleting it.
- Update `task.md` with the Rationale.

3. **Pattern Matching & Performance Budget**:

- **CHECK**: Scan `coding-standards.md` and `docs/adr/` for past mistakes.
- **BUDGET**: Estimate the latency impact. If it adds >100ms to the critical path, you MUST propose a simpler alternative.
- **CHALLENGE**: If the proposed idea resembles a past mistake or violates a "Strict Lesson," surface this immediately.
- **PROMPT**: If a conflict exists, suggest: "This looks like [Past Mistake/Violation]. Before we proceed, we should run `/learn` to update our constraints"

4. **Delete the Part or Process**:

- Actively propose removing at least one component, feature, or abstraction.
- **Reject Interface Bloat**: Push for direct access via a unified `*AppEnv` context. Update `task.md` with what was cut.

5. **Observability & Disk-Parity Strategy**:

- **Observability**: Define the "Success Metric" (e.g., search completion rate). Identify required logs/metrics.
- **Parity**: Plan for a file-based SQLite audit in Phase 3 to catch locking issues.

6. **Task Initialization**:

- Populate the `#Execution Plan` section with `[ ]` items derived from the Decision Log for Phase 2 and Phase 3.

## Phase 2: Autonomous Execution Loop (TDD)

**Pre-condition**: Read .agents/skills/go-tdd/SKILL.md before starting this phase.

1. **RED**: Write failing tests. Run `go test`.
2. **GREEN**: Write implementation (including logs/metrics). 
- Run `go test`. 
- If you fail to achieve GREEN after 3 attempts, HALT. Read the raw tracebacks to the user, hypothesize the flaw in the test or implementation, and WAIT for user guidance.
3. REFACTOR: Run `go run ./cmd/verify critique`. Run `go run ./cmd/verify heal`. Fix any remaining checks.
4. **COMMIT**: Once tests and lint pass, execute `git commit -m "feat(<scope>): implement <idea>"` natively.

## Phase 3: Audit & Resilience

1. **Security Audit**: Run your local SAST/security tooling (e.g., `gosec ./...`). Fix any P0/P1/P2 defects immediately and amend the commit.
2. **Performance Audit**: Verify query execution plans using `EXPLAIN QUERY PLAN`.
3. **Disk-Parity Test**: Run tests against a temporary file-backed SQLite DB (not just `:memory:`) to verify WAL-mode/concurrency behavior.
4. **Chaos/Resilience**: Identify the weakest architectural boundary of the new feature based on the project structure:

- `repository/`: Inject database timeouts or connection drops.
- `handler/` & `middleware/`: Inject malformed payloads or test rate limit exhaustion.
- `service/`: Force business logic edge cases or unexpected interface nil returns.

5. **The Resilience Halt**: If chaos tests break the app, HALT and explain. Propose 2-3 patterns (e.g., Circuit Breaker). WAIT for user decision.
6. **Contract Verification**:

- Run `go run cmd/verify/main.go template-drift`
- Run `go run cmd/verify/main.go api-spec`

6. **UI Verification**:
   - **Pre-condition**: Read .agents/skills/browser-verify/SKILL.md before any browser task.
   - use `browser_subagent` to capture a screenshot of the "Find Food" flow and embed it in the walkthrough.
- Verify the "Final Truth" against the persona requirements (e.g., "Is the pivot visible?").
- **HARD GATE**: You are forbidden from drafting the `walkthrough.md` or summarizing completion until this step is documented with a screenshot and checked off in `task.md`.

**Completion**: When all phases are complete and the final commit is made, summarize the architectural decisions and test coverage for the user in a single, concise chat message.