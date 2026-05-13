# ADR [006]: Quota Automation Gates
**Date**: 2026-05-13 **Status**: Proposed

## 1. Context & User Problem
The current "QUOTA PROTECTION GATE" relies on markdown instructions that models can ignore. High-tier models (Pro/Opus) are often used for iterative product execution (TDD loops), which is extremely expensive and slow compared to Flash. Additionally, "context rot" occurs when preflight rules and history balloon, degrading agentic accuracy and increasing costs.

## 2. Decision
We will implement deterministic CLI gates in `cmd/verify` to:
1. **Quota Gate**: Verify that any commit indicating high-tier model usage (detected via metadata or reasoning markers) must contain the word `OVERRIDE` if it mutates product code.
2. **Preflight Tax**: Enforce a strict byte-size limit on the "bootstrap" bundle (`AGENTS.md`, `RESOLVER.md`, `coding-standards.md`, `verify-manifest.yaml`).
3. **Command Integration**: Wire these into the `precommit` gate to prevent expensive mistakes from reaching the remote repository.

## 3. The Complexity Kill-Switch (Rationale)
* **User Value**: Drastically reduces the cost of autonomous development loops and prevents "context rot" from slowing down feature delivery.
* **Performance Budget**: < 50ms impact on `precommit` (Deterministic file-size and string checks).
* **Minimalism Check**: Replaces manual "Context Density" monitoring with automated, deterministic enforcement.

## 4. Consequences
* **Technical Tradeoffs**: Requires manual intervention/overrides during major architectural changes that genuinely require high-tier reasoning.
* **Observability**: Failures will be explicitly flagged in the terminal and remote CI logs.
* **SQLite Impact**: None.

## 5. Alternatives Considered
* **Dynamic Token Counting**: Rejected as too complex (requires external dependencies or API calls). Byte-size is a reliable, conservative proxy for token usage.
