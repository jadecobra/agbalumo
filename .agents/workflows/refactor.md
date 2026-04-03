---
description: Workflow for refactoring internal components without changing external API or CLI contracts.
---

# Refactor Workflow

> **Rule**: Refactoring must NOT change the external API or CLI contracts. Logic changes must be verified by existing tests and mandatory drift checks. We do not write code that breaks. We are 10x.

---

## 0. Initialize Workflow State

Before any code is modified, initialize the state machine to track gate progress.

// turbo
1. Run `./scripts/agent-exec.sh init <refactor-description> refactor`
// turbo
2. Run `./scripts/agent-exec.sh set-phase RED`

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
  - **PASS**: Existing tests pass, or new "safety" tests fail as expected.
  - **Note**: In a pure refactor, you may manually mark this as PASS if existing tests are sufficient.
  - `DEBUG`: `./scripts/agent-exec.sh verify red-test`

### 1b. Verify Contract Stability

> **Persona: SystemsArchitect** — Ensure the refactor doesn't accidentally change the API or CLI (One-Way Door).

- **Gate: `api-spec`**
  - **PASS**: `./scripts/agent-exec.sh verify api-spec` passes.
  - **FAIL**: Drift detected. Revert contract change or justify to the SystemsArchitect and ProductOwner.

---

## 2. Refactoring Implementation (GREEN & REFACTOR)

Modify the code while keeping tests green.

### 2a. Implement Refactor (GREEN)

> **Persona: Backend** — Clean up code, improve performance, or fix technical debt.

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
  task pre-commit
  ```
- **Gate: `lint`**
- **Gate: `coverage`** (Coverage MUST NOT drop).
- **Advisory: `context-cost`** (Target TokenRMS < 110, monitor only).

---

## 3. Final Verification

### 3a. UI & Browser (If applicable)
- *Run*:
  ```bash
  ./scripts/verify_restart.sh
  ```
@[.agents/rules/browser-url.md]
- **Gate: `browser-verification`**

### 3b. Chaos Contract Audit
> **Persona: ChaosMonkey** — Attempt to silently modify a contract (API/CLI) without updating documentation. Verify that the `api-spec` gate correctly identifies and blocks the change.

---

## 6. Final Reset

// turbo
1. Run `./scripts/agent-exec.sh set-phase IDLE`
// turbo
2. Run `./scripts/agent-exec.sh init none`

---

## Completion Checklist

- [ ] **Gate: `red-test`** - Safety tests verified.
- [ ] **Gate: `api-spec`** - Drift checks PASSED (Contracts preserved).
- [ ] **Gate: `implementation`** - Refactor complete and tests pass.
- [ ] `go test -race ./...` passes.
- [ ] **Gate: `lint`** - `task pre-commit` passed.
- [ ] **Gate: `coverage`** - No coverage drop.
- [ ] **Context-Cost Check** - TokenRMS awareness verified.
- [ ] **Gate: `chaos-verify`** - Refactor survived fault injection.
- [ ] `task.md` updated.
- [ ] **AUTO-COMMIT**: Execute git commit automatically with a short, imperative message (e.g., "Refactor: extract helper function"). DO NOT wait for the user to explicitly tell you to commit.
- [ ] **FINALIZE**: Instruct the agent to use the `mcp_mcp-memory-service_memory_store` tool to save the refactor completion and context.
