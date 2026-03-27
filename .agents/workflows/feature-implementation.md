---
description: Workflow for implementing new features with strict TDD, API-first design, and full validation (CLI, UI, Browser)
---

# Feature Implementation & Validation Workflow

> **Rule**: A feature is NOT done until it passes all of these layers. We do not write code without tests. We are 10x.

---

## 0. Initialize Workflow State

@[.agent/workflows/feature-implementation/planning.md]

---

@[.agent/workflows/feature-implementation/red.md]

---

## 2. Core Implementation (GREEN & REFACTOR)

Write the minimum code necessary to satisfy the API spec and pass the test.

@[.agent/workflows/feature-implementation/green.md]

@[.agent/workflows/feature-implementation/refactor.md]

---

## 3. Command Line & Integration Testing

Validate that the application can be used from the command line/terminal as an admin or user.

### 3a. CLI Test

> **Persona: SDET** — Write table-driven tests. Cover end-to-end CLI workflows. Validate both success and error paths. Never write production code in this step.

- Test with the `agbalumo` CLI.
- Ensure the newly implemented feature can be invoked from the command line and works end-to-end as intended.
- *Run*:
  ```bash
  // turbo
  go test -v ./cmd/...
  ```
- **Gate: `implementation` (CLI)**
  - **PASS**: CLI command returns expected success output or `go test ./cmd/...` passes.
  - **FAIL**: CLI returns error or unexpected data.

---

## 4. UI & Browser Verification

Ensure the implemented feature is fully functional and pixel-perfect in the UI.

### 4a. Build & Restart Server
- *Run*:
  ```bash
  ./scripts/verify_restart.sh
  ```
- **Gate**: Server MUST compile successfully and remain running on expected ports.

### 4b. Programmatic UI Test

> **Persona: SDET + UI/UX** — Write integration/UI tests simulating browser interactions. Verify HTMX renders, form posts, and admin workflows. Ensure pixel-perfect alignment, consistent spacing, and brand adherence (Orange `#FF5E0E`, Green `#2D5A27`).

- Add or update integration/UI tests (e.g. HTMX renders, form posts) calling the handlers to simulate browser interactions programmatically.
- **Gate**: Programmatic UI tests pass and accurately reflect user/admin workflows.

### 4c. Browser Subagent Verification

> **Persona: UI/UX + QA** — Verify premium feel: micro-interactions, smooth transitions, responsive layout. Ensure FCP < 1.0s. Capture browser recordings to verify responsiveness. Check agbalumo brand consistency.

- Use the `browser_subagent` tool with a detailed task:
  ```
  Task: "Navigate to the URL defined in .agents/rules/browser-url.md. Check [FEATURE].
    1. Act as user/admin. Navigate to the view.
    2. Attempt to use feature (fill forms, click buttons).
    3. Verify success states and error states.
    4. Perform basic accessibility check (tab order, visual focus)."
  RecordingName: "verify_feature_name_ui"
  ```
- **Gate: `browser-verification`**
  - **PASS**: Browser subagent report confirms feature works as intended; visual audit passes; zero console errors.
  - **FAIL**: Interactions fail, data missing, or visual breakage detected.

---

## 6. Final Reset

After all verification is complete, reset the workflow state and restart the server:
// turbo
1. Run `./scripts/agent-exec.sh workflow set-phase IDLE`
// turbo
2. Run `./scripts/agent-exec.sh workflow init none`

Then restart the server:
@[.agent/workflows/restart-server.md]

---

## Completion Checklist

A feature is **DONE** when ALL boxes are checked:

- [ ] **Gate: `red-test`** - Unit test was written FIRST and **failed** (Red).
- [ ] **Gate: `api-spec`** - Requirements defined in `@[docs/api.md]` and `@[docs/openapi.yaml]`.
- [ ] **Gate: `implementation`** - Logic implemented and tests **passed** (Green).
- [ ] **Gate: `implementation` (CLI)** - Verified via `agbalumo` CLI as user/admin.
- [ ] Programmatic UI tests as user/admin verify it can be used from the UI.
- [ ] `go test -race ./...` passes (no regressions).
- [ ] **Gate: `lint`** - `task pre-commit` linting passed.
- [ ] **Gate: `coverage`** - `task pre-commit` coverage threshold met.
- [ ] `./scripts/verify_restart.sh` executed successfully (server running).
- [ ] **Gate: `browser-verification`** - Browser subagent verified feature natives UI.
- [ ] Recording artifact saved with descriptive name.
- [ ] `task.md` updated with completed status.
- [ ] `spec.md` reviewed and updated if needed.
- [ ] `@[.agent/workflows/restart-server.md]` was run after verification.
- [ ] **AUTO-COMMIT**: Execute git commit automatically with a short, imperative message. DO NOT wait for the user to explicitly tell you to commit.
- [ ] **FINALIZE**: Instruct the agent to use the `mcp_mcp-memory-service_memory_store` tool to save the feature completion and context.