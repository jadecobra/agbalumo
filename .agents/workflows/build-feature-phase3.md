---
description: Phase 3: Audit & Resilience
---
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
7. **UI Verification**:
   - **Pre-condition**: Read .agents/skills/browser-verify/SKILL.md before any browser task.
   - use `browser_subagent` to capture a screenshot of the "Find Food" flow and embed it in the walkthrough.
- Verify the "Final Truth" against the persona requirements (e.g., "Is the pivot visible?").
- **HARD GATE**: You are forbidden from drafting the `walkthrough.md` or summarizing completion until this step is documented with a screenshot and checked off in `task.md`.
8. **Knowledge Extraction (Skill & Tool Audit)**:
   - Review the session's `task.md` Decision Log and git log.
   - If any multi-step procedure was repeated ≥2 times during this session, or required ≥2 correction attempts: extract it into a new Skill in `.agents/skills/`.
   - If any manual check could be automated with a deterministic pass/fail: propose a new `verify` subcommand.
   - Update `.agents/verify-manifest.yaml` and `AGENTS.md` with any new skills or tools.
**Completion**: When all phases are complete and the final commit is made, summarize the architectural decisions and test coverage for the user in a single, concise chat message.
