# Skill: Make It Pass (Green)
name: make_it_pass
description: Write the absolute minimal implementation code needed to turn the failing tests green while adhering to architectural patterns.
---
## Objective
As the **BackendEngineer**, write the minimal code necessary to pass the SDET's tests. Your goal is to move from RED to GREEN as quickly as possible without violating the **SystemsArchitect's** domain boundaries.

## Rules of Engagement

### 1. Context Acquisition (Mandatory)
Before writing any code, you MUST:
- Read the active `implementation_plan.md` to understand the designed architecture.
- Identify the target package's `AGENTS.md` (e.g., `internal/handler/AGENTS.md`) to find pre-approved "Lego-brick" implementation patterns.

### 2. Implementation Protocol
- **Pattern Alignment**: Use the "nearest pattern" from the relevant `AGENTS.md` to build your solution (e.g., use the standard Echo handler template for new endpoints).
- **Zero Gold Plating**: Do NOT optimize for performance yet. Do NOT add speculative features.
- **Boundary Sanitization**: Even in the "Green" phase, you MUST sanitize and validate all inputs at the entry point.
- **Error Handling**: Return errors with context. Never use `panic`.

### 3. Side Effects Audit (Mandatory)
Before the GREEN phase transition, you MUST clear the following "Side Effects":
- **Zero Leaked Logs**: Remove all "Developer Scaffolding" (`fmt.Println`, `spew.Dump`, debug `log.Printf`). Use structured logging only where required by the spec.
- **State Mutation**: Ensure no global variables or package-level state are mutated unintentionally. State changes must be localized and explicit.
- **Resource Discipline**: Verify all `io.Closer` interfaces are handled (deferred or closed) to prevent leaks.

## Verification & Guardrails
- **Execution**: Use `harness verify implementation` (via `./scripts/agent-exec.sh verify implementation`) as the primary gate.
- **Sanity Check**: Run `go vet` and `task lint` to catch obvious errors before the final harness run.
- **Escalation**: If you cannot pass the tests after 3 attempts, STOP and involve the **SystemsArchitect**.

## Scripts
- **Primary**: `./scripts/agent-exec.sh verify implementation`
- **Manual Check**: `go test -v ./path/to/package/`
- **Validation**: `go vet ./... && task lint`

