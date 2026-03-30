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

### 3. Verification & Guardrails
- **Execution**: Use `./scripts/agent-exec.sh verify implementation` to run the tests.
- **Sanity Check**: Run `go vet` to catch obvious errors before verifying.
- **Escalation**: If you cannot pass the tests after 3 attempts, or if the requirements conflict with the existing patterns, STOP and involve the **SystemsArchitect**.

## Scripts
- **Primary**: `./scripts/agent-exec.sh verify implementation`
- **Manual Check**: `go test -v ./path/to/package/`
- **Validation**: `go vet ./...`
