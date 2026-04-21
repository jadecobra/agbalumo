# ADR [006]: Maintenance Audit Refactoring (Phase 2)
**Date**: 2026-04-21 **Status**: Accepted

## Context
Following the successful decomposition of the `verify` CLI (Phase 1), `internal/maintenance/audit.go` has been identified as the next high-density bottleneck. It currently manages disparate logic including low-level security grep-checks, complex golangci-lint reporting (ChiefCritic), and automated healing logic, leading to context bloat.

## Decision
We will decompose `internal/maintenance/audit.go` into domain-specific files:
- `internal/maintenance/security.go`: Standalone security validations (vet, headers, fly.toml, vuln, XSS).
- `internal/maintenance/chief_critic.go`: Golangci-lint output parsing and Agent-Native summary reporting.
- `internal/maintenance/heal.go`: Automated structural corrections (fieldalignment).
- `internal/maintenance/util.go`: Shared execution utilities like `runTool`.

`audit.go` will be stripped of implementations and focus solely on high-level orchestration of security checks.

## Consequences
- **Context Precision**: AI Agents can now target security checks or formatting logic without loading the entire audit suite.
- **Improved Signal**: Separating reporting logic (ChiefCritic) from validation logic (Security) makes the codebase easier to evolve independently.
- **Contract Stability**: External maintenance interfaces remain unchanged; only internal file structure is optimized.
