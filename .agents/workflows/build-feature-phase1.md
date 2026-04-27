---
description: Phase 1: Product Architecture & Planning
---
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
