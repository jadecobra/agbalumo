# Quota Automation & Performance Tooling

## Decision Log
- **Product Interrogation (Kill-Switch)**: This feature does not directly impact the 60-second end-user search latency. However, it strictly enforces token limits and context boundaries, preventing AI execution from ballooning the development loop latency. Value = Dev-velocity & Cost protection.
- **Pattern Match Against History**: This codifies the "QUOTA PROTECTION GATE" defined in `AGENTS.md` and expands upon ADR 005 (Token Density Refactoring).
- **Architectural Decision (ADR)**: We will create a new ADR `2026-05-13-quota-automation-gates.md` documenting the enforcement of deterministic CLI gates for model tiering and context size limits.
- **Delete the Part**: Instead of a complex dynamic LLM proxy or AST-based minifier, we will use simple, deterministic CLI tooling (`cmd/verify/*`) to enforce strict size thresholds and commit message flags.

## Execution Plan

### Phase 2: Autonomous Execution Loop (TDD)
- [x] Create `docs/adr/2026-05-13-quota-automation-gates.md`.
- [x] Implement `internal/maintenance/quota.go` with tests to verify commit messages for High-Tier usage violations (unless `OVERRIDE` is present) and measure preflight file sizes.
- [x] Implement `cmd/verify/quota_gate.go` to expose the CLI command `quota-gate`.
- [x] Implement `cmd/verify/preflight_tax.go` to expose the CLI command `preflight-tax` (fails if `AGENTS.md`, `RESOLVER.md`, `coding-standards.md`, and `verify-manifest.yaml` combined exceed 15KB).
- [x] Wire `quota-gate` and `preflight-tax` into the standard `precommit` runner located in `cmd/verify/ci.go`.
- [x] Update `.agents/verify-manifest.yaml` to register the new commands.

### Phase 3: Audit & Resilience
- [ ] Run `go run ./cmd/verify critique --baseline=HEAD~1`.
- [ ] Verify execution passes locally: `go run ./cmd/verify precommit`.
- [ ] Execute `./scripts/pushw.sh` and ensure the remote CI is green.
