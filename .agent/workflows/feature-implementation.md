---
description: Workflow for implementing new features with strict TDD, API-first design, and full validation (CLI, UI, Browser)
---

# Feature Implementation & Validation Workflow

> **Rule**: A feature is NOT done until it passes all of these layers. We do not write code without tests. We are 10x.

---

## 1. API Specification & Red Test (Design First)

Before writing any implementation code, prove the feature doesn't exist and define its contract.

### 1a. Write Failing Unit Test (RED)
- *Role*: SDET Agent
- *File*: Create or update the relevant `*_test.go` file.
- Write a test expecting the new feature to work (e.g. hitting the hypothetical endpoint).
- *Run*:
  ```bash
  // turbo
  go test -v -run TestNewFeatureName ./internal/package_name/...
  ```
- **Gate**: The test MUST **FAIL** (e.g., 404 Not Found, or missing struct field). It proves the implementation is not in the API yet.

### 1b. Update API Specifications (Source of Truth)
- *Role*: Lead Architect
- *Files*: Update `@[docs/api.md]` and `@[docs/openapi.yaml]` to map out the exact request/response/path for the feature. This acts as the absolute source of truth.
- Include appropriate validation rules and security requirements.
- *Lint Spec* (Optional but recommended):
  ```bash
  swagger-cli validate docs/openapi.yaml
  ```
- **Gate**: Specification accurately describes the goal.

---

## 2. Core Implementation (GREEN & REFACTOR)

Write the minimum code necessary to satisfy the API spec and pass the test.

### 2a. Implement Logic (GREEN)
- *Role*: Backend Agent
- Write the logic in `handler/`, `service/`, etc.
- Also ensure tests cover `400 Bad Request` schema validations and `401/403` auth checks defined in the spec.
- *Run*:
  ```bash
  // turbo
  go test -v -run TestNewFeatureName ./internal/package_name/...
  ```
- **Gate**: Test MUST **PASS**. The implementation completely meets the specification in `@[docs/api.md]` and `@[docs/openapi.yaml]`.

### 2b. Refactor & Regression
- *Run full regression*:
  ```bash
  // turbo
  go test -race ./...
  ```
- **Gate**: Zero regressions.

### 2c. Pre-Commit Quality Gate
- *Run*:
  ```bash
  ./scripts/pre-commit.sh
  ```
- **Gate**: Script MUST exit 0 (Coverage >= 82.5%, no lint errors).

---

## 3. Command Line & Integration Testing

Validate that the application can be used from the command line/terminal as an admin or user.

### 3a. CLI Test
- *Role*: SDET Agent
- Test with the `agbalumo` CLI.
- Ensure the newly implemented feature can be invoked from the command line and works end-to-end as intended.
- *Run*:
  ```bash
  // turbo
  go test -v ./cmd/...
  ```
- **Gate**: The feature can be successfully invoked and verified via CLI for user/admin personas.

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
- *Role*: SDET / UI Engineer
- Add or update integration/UI tests (e.g. HTMX renders, form posts) calling the handlers to simulate browser interactions programmatically.
- **Gate**: Programmatic UI tests pass and accurately reflect user/admin workflows.

### 4c. Browser Subagent Verification
- *Role*: QA / UI Engineer
- Use the `browser_subagent` tool with a detailed task:
  ```
  Task: "Navigate to https://localhost:8443. Check [FEATURE].
    1. Act as user/admin. Navigate to the view.
    2. Attempt to use feature (fill forms, click buttons).
    3. Verify success states and error states.
    4. Perform basic accessibility check (tab order, visual focus)."
  RecordingName: "verify_feature_name_ui"
  ```
- **Gate**:
  - Feature exists and data renders correctly.
  - Visual quality is premium (brand colors, layout).
  - No JavaScript or CSP errors in the console.

---

## 5. Final Reset

After every 3-step verification is complete, reset the application state by running the restart server workflow:
@[.agent/workflows/restart-server.md]

---

## Completion Checklist

A feature is **DONE** when ALL boxes are checked:

- [ ] Unit test was written FIRST and **failed** (Red), verifying implementation is not in the API
- [ ] Requirements defined cleanly in `@[docs/api.md]` and `@[docs/openapi.yaml]` as the source of truth
- [ ] Minimal code implemented, edge cases (400, 401, 403) covered, and tests **passed** (Green), verifying it meets specification
- [ ] Tested with `agbalumo` CLI as user/admin to verify it works from the command line
- [ ] Programmatic UI tests as user/admin verify it can be used from the UI
- [ ] `go test -race ./...` passes (no regressions)
- [ ] `./scripts/pre-commit.sh` exits 0 (strict coverage maintained)
- [ ] `./scripts/verify_restart.sh` executed successfully (server running)
- [ ] Browser subagent as user/admin verified the feature works from the UI natively
- [ ] Recording artifact saved with descriptive name
- [ ] `task.md` updated with completed status
- [ ] `spec.md` reviewed and updated if needed
- [ ] `@[.agent/workflows/restart-server.md]` was run after verification
- [ ] Commit with short, imperative message