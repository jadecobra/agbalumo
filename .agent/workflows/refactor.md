---
description: Workflow for refactoring internal components without changing external API or CLI contracts.
---

# Refactor Workflow

> **Rule**: Refactoring must NOT change the external API or CLI contracts. Logic changes must be verified by existing tests and mandatory drift checks. We do not write code that breaks. We are 10x.

---

## 0. Initialize Workflow State

Before any code is modified, initialize the state machine to track gate progress.

// turbo
1. Run `./scripts/agent-exec.sh workflow init <refactor-description> refactor`
// turbo
2. Run `./scripts/agent-exec.sh workflow set-phase RED`

---

## 1. Safety Checks & Baseline

Define the scope and ensure the current codebase is stable.

### 1a. Identify Regressions (RED)

> **Persona: SDET** — Identify existing tests that cover the area being refactored. If coverage is low, write additional tests FIRST to capture current behavior.

- *Run*:
  ```bash
  // turbo
  go test -v ./internal/package_name/...
  ```
- **Gate: `red-test`**
  - **PASS**: Existing tests pass, or new "safety" tests fail as expected (if written to prove a bug or missing coverage).
  - **Note**: In a pure refactor, you might manually mark this as PASS if existing tests are sufficient.
  - `DEBUG`: `./scripts/agent-gate.sh red-test`

### 1b. Verify Contract Stability

> **Persona: Lead Architect** — Ensure the refactor doesn't accidentally change the API or CLI.

- **Gate: `api-spec`**
  - **PASS**: `./scripts/agent-gate.sh api-spec` passes (no drift in `docs/api.md` or CLI contracts).
  - **FAIL**: Drift detected. You MUST either revert the contract change or justify why it's necessary (which might turn this into a `feature` workflow).

---

## 2. Refactoring Implementation (GREEN & REFACTOR)

Modify the code while keeping tests green.

### 2a. Implement Refactor (GREEN)

> **Persona: Backend** — Clean up code, improve performance, or fix technical debt. Maintain functional equivalence.

- *Run*:
  ```bash
  // turbo
  go test -v ./internal/package_name/...
  ```
- **Gate: `implementation`**
  - **PASS**: All tests pass.
  - **FAIL**: Regressions found.

### 2b. Regression Suite
- *Run full regression*:
  ```bash
  // turbo
  go test -race ./...
  ```

### 2c. Quality Gate
- *Run*:
  ```bash
  ./scripts/pre-commit.sh
  ```
- **Gate: `lint`**
- **Gate: `coverage`** (Coverage MUST NOT drop).

---

## 3. Final Verification

### 3a. UI & Browser (If applicable)
- *Run*:
  ```bash
  ./scripts/verify_restart.sh
  ```
- **Gate: `browser-verification`** (Optional, but recommended if UI logic was touched).

---

## 6. Final Reset

// turbo
1. Run `./scripts/agent-exec.sh workflow set-phase IDLE`
// turbo
2. Run `./scripts/agent-exec.sh workflow init none`

---

## Completion Checklist

- [ ] **Gate: `red-test`** - Safety tests verified.
- [ ] **Gate: `api-spec`** - Drift checks PASSED (Contracts preserved).
- [ ] **Gate: `implementation`** - Refactor complete and tests pass.
- [ ] `go test -race ./...` passes.
- [ ] **Gate: `lint`** - `./scripts/pre-commit.sh` passed.
- [ ] **Gate: `coverage`** - No coverage drop.
- [ ] `task.md` updated.
- [ ] Commit with short, imperative message (e.g., "Refactor: extract helper function").
- [ ] **FINALIZE**: Instruct the agent to use the `mcp_mcp-memory-service_memory_store` tool to save the refactor completion and context.
