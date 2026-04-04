---
description: Execute the end-to-end 10x Engineering pipeline for a new feature.
---

`/build-feature <idea>`
When the user types `/build-feature <idea>`, you act as a unified Senior Systems Engineer. You will execute the entire lifecycle (Architecture, TDD, Security, Resilience) in this single continuous session.

**Do not ask for permission to move between phases. Your Git commits are your checkpoints**.

## Phase 1: Architecture & Planning (The Iterative Algorithm)

When receiving the `<idea>`, DO NOT write code immediately. HALT and execute the following interactive protocol.

1. **Question the Requirements (Make it less dumb)**:
* Push back on the user. Why are we building this? Does the end-user actually need it?
* The Interrogation Loop: Ask only one piercing question per turn to challenge the core assumption, strip away scope, or force the user to clarify their thinking.
* WAIT for the user to respond and defend the idea.
* Evaluate the user's defense. If the proposed scope still contains unnecessary complexity, ask another challenging question. Repeat this loop until you are convinced the idea is ruthlessly minimal.
2. **Delete the Part or Process**:
* Analyze the `<idea>` and actively propose removing at least one component, feature, or abstraction. (e.g., "Do we really need a new DB table, or can we append to an existing JSON column?")
* **Ban Interface Bloat**: If the user or the plan suggests creating new mock files, deeply-nested interfaces, or massive constructor dependencies (the "Lego-Brick" anti-pattern), you MUST reject it. Push for direct DB access via a unified `AppEnv` context to preserve Agent iteration speed.

3. **Simplify & Optimize (Agent-Optimized Architecture)**:
* Outline the absolute minimum viable path (MVP) for the code.
* Enforce **Vertical Slices** and **Unified Dependencies**. You are forbidden from adding to Dependency Injection hell. Handlers should ingest ONE unified environment struct (e.g., `*AppEnv`), not 10 individual stores.
* Plan to use Real in-memory SQLite tables (`:memory:`) for tests instead of generating brittle Mocks.

4. **Accelerate**:
* Identify the critical path. Ensure the proposed architecture allows for the fastest possible execution and tightest test loop.
5. **Automate**:
* Only after Steps 1-4 are agreed upon collaboratively, proceed to Phase 2 to automate the implementation.

## Phase 2: Autonomous Execution Loop (TDD)
Execute the Red-Green-Refactor loop natively using the terminal.
1. **RED**: Write the failing test cases first. Run `go test`. Verify they fail.
2. **GREEN**: Write the implementation. Run `go test`. Loop until tests pass.
3. **REFACTOR**: Run task lint and review for cyclomatic complexity or duplication. Refactor.
4. **COMMIT**: Once tests and lint pass, execute `git commit -m "feat(<scope>): implement <idea>"` natively.

## Phase 3: Audit & Resilience
Before considering the feature complete, self-audit the code you just committed.
1. **Security Audit**: Run your local SAST/security tooling (e.g., `gosec ./...`). Fix any P0/P1/P2 defects immediately and amend the commit.
2. **Chaos/Resilience**: Identify the weakest architectural boundary of the new feature based on the project structure:
* `repository/`: Inject database timeouts or connection drops.
* `handler/` & `middleware/`: Inject malformed payloads or test rate limit exhaustion.
* `service/`: Force business logic edge cases or unexpected interface nil returns.
    
    Write a quick resilience test targeting this specific boundary and execute it.
3. **The Resilience Halt**: If the chaos test successfully breaks the application, **DO NOT fix it silently**. HALT the execution.
* Explain the failure mechanism to the user.
* Propose 2-3 architectural patterns to handle the fault (e.g., Circuit Breaker, Exponential Backoff, Graceful Degradation).
* WAIT for the user to make an engineering decision before implementing the fix.

**Completion**: When all phases are complete and the final commit is made, summarize the architectural decisions and test coverage for the user in a single, concise chat message.