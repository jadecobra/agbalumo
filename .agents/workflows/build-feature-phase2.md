---
description: Phase 2: Autonomous Execution Loop (TDD)
---
## Phase 2: Autonomous Execution Loop (TDD)
**Pre-condition**: Read .agents/skills/go-tdd/SKILL.md before starting this phase.
1. **RED**: Write failing tests. Run `go test`.
2. **GREEN**: Write implementation (including logs/metrics). 
- Run `go test`. 
- If you fail to achieve GREEN after 3 attempts, HALT. Read the raw tracebacks to the user, hypothesize the flaw in the test or implementation, and WAIT for user guidance.
3. REFACTOR: Run `go run ./cmd/verify critique`. Run `go run ./cmd/verify heal`. Fix any remaining checks.
4. **COMMIT**: Once tests and lint pass, execute `git commit -m "feat(<scope>): implement <idea>"` natively.
