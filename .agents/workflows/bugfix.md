---
description: Workflow for fixing bugs with reproduction tests and contract validation.
---

# Bugfix Workflow

> **Rule**: No bug is fixed until a reproduction test exists and passes. Fixes must NOT break existing API/CLI contracts unless the bug WAS the contract. We do not write code that breaks. We are 10x.

---

## 0. Initialize Workflow State

Before any code is modified, initialize the state machine.

// turbo
1. Run `./scripts/agent-exec.sh workflow init <bug-id-or-description> bugfix`
// turbo
2. Run `./scripts/agent-exec.sh workflow set-phase RED`

---

## 1. Reproduction (RED)

Prove the bug exists.

### 1a. Write Reproduction Test (RED)

> **Persona: SDET** — Write a test that specifically triggers the reported bug. The test MUST fail.

- *File*: Update the relevant `*_test.go` file.
- *Run*:
  ```bash
  // turbo
  go test -v -run TestBugReproduction ./internal/package_name/...
  ```
- **Gate: `red-test`**
  - **PASS**: Test fails exactly as described in the bug report.
  - **FAIL**: Test passes (bug not reproduced) or fails for wrong reason.
  - `DEBUG`: `./scripts/agent-gate.sh red-test "Error message pattern"`

### 1b. Verify Contract Stability

> **Persona: Lead Architect** — Ensure the fix doesn't break external contracts.

- **Gate: `api-spec`**
  - **PASS**: `./scripts/agent-gate.sh api-spec` passes.
  - **FAIL**: Fix requires contract change (must be justified).

---

## 2. Fix Implementation (GREEN)

Implement the minimal fix.

### 2a. Implementation (GREEN)

> **Persona: Backend** — Fix the bug with the minimal necessary code.

- *Run*:
  ```bash
  // turbo
  go test -v -run TestBugReproduction ./internal/package_name/...
  ```
- **Gate: `implementation`**
  - **PASS**: Reproduction test passes.
  - **FAIL**: Test still fails.

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
- **Gate: `coverage`**

---

## 3. Browser/UI Verification (If applicable)

- *Run*:
  ```bash
  ./scripts/verify_restart.sh
  ```
@[.agents/rules/browser-url.md]
- **Gate: `browser-verification`**

---

## 6. Final Reset

// turbo
1. Run `./scripts/agent-exec.sh workflow set-phase IDLE`
// turbo
2. Run `./scripts/agent-exec.sh workflow init none`

---

## Completion Checklist

- [ ] **Gate: `red-test`** - Bug reproduced with failing test.
- [ ] **Gate: `api-spec`** - Drift checks PASSED.
- [ ] **Gate: `implementation`** - Bug fixed and tests pass.
- [ ] `go test -race ./...` passes.
- [ ] **Gate: `lint`** - `./scripts/pre-commit.sh` passed.
- [ ] **Gate: `coverage`** - Threshold met.
- [ ] `task.md` updated.
- [ ] **AUTO-COMMIT**: Execute git commit automatically with a short, imperative message (e.g., "Fix: resolve nil pointer in listing handler"). DO NOT wait for the user to explicitly tell you to commit.
- [ ] **FINALIZE**: Instruct the agent to use the `mcp_mcp-memory-service_memory_store` tool to save the bugfix completion and context.
