---
description: Workflow for implementing new features with strict TDD and 3-layer verification
---

# Feature Implementation Workflow (3-Layer Verification)

> **Rule**: A feature is NOT done until it passes all 3 layers.

---

## Layer 1: Unit Tests (Red → Green → Refactor)

### 1a. Write Failing Test (RED)

- *Role*: SDET Agent
- *File*: Create or update the relevant `*_test.go` file.
- *Run*:
  ```bash
  // turbo
  go test -v -run TestNewFeatureName ./internal/package_name/...
  ```
- **Gate**: The test MUST **FAIL**. If it passes, the test is not testing new behavior. Rewrite it.
- *Verify*: Failure is due to *missing logic*, not a compilation error (unless adding a new API/struct).

### 1b. Implement (GREEN)

- *Role*: Backend Agent
- *Constraint*: Write the **minimal** code to make the test pass. Do not add untested features.
- *Run*:
  ```bash
  // turbo
  go test -v -run TestNewFeatureName ./internal/package_name/...
  ```
- **Gate**: The test MUST **PASS**.

### 1c. Refactor

- *Role*: Backend Agent / Lead Architect
- *Run full regression*:
  ```bash
  // turbo
  go test -race ./...
  ```
- **Gate**: ALL existing tests MUST still pass. Zero regressions.

---

## Layer 2: CLI / Integration Test

Validates the feature works end-to-end through the application's command layer.

### 2a. Run Pre-Commit Quality Gate

```bash
./scripts/pre-commit.sh
```

This enforces:
- `gofmt` formatting
- `go mod tidy`
- `go vet`
- Race detection
- Coverage >= 87.8% threshold (NEVER lower this — write more tests instead)
- Secret scanning

**Gate**: Script MUST exit 0.

### 2b. Build & Restart Server

```bash
./scripts/verify_restart.sh
```

This enforces:
- Runs pre-commit checks (Layer 2a)
- Builds CSS assets (`npm run build:css`)
- Compiles binary (`go build -o bin/agbalumo main.go`)
- Kills old process on :8443/:8080
- Starts new server
- Confirms server stays alive for 2s

**Gate**: Server MUST be running after script completes.

### 2c. CLI Smoke Test (if applicable)

For command-line features (seed, serve config, etc.):

```bash
// turbo
go test -v ./cmd/...
```

**Gate**: All cmd tests pass.

---

## Layer 3: Browser Subagent Verification

Validates the feature works for a real user in the browser.

### 3a. Launch Browser Subagent

Use the `browser_subagent` tool with a task like:

```
Task: "Navigate to https://localhost:8443. Verify [SPECIFIC FEATURE].
  1. [Step to reach the feature, e.g. 'Click on Create Listing']
  2. [Step to interact, e.g. 'Fill in Title with Senior Go Engineer']
  3. [Step to verify, e.g. 'Confirm the listing card shows Company Name']
  Return: Screenshot proof and pass/fail for each step."

RecordingName: "verify_feature_name"
```

### 3b. What to Verify

| Check | How |
|:---|:---|
| **Feature exists** | Element is visible and interactive |
| **Data renders** | Correct values appear (not empty, not placeholder) |
| **Visual quality** | Spacing, colors, typography match brand (Orange #FF5E0E, Green #2D5A27) |
| **Responsiveness** | Resize to mobile (375px) and verify layout |
| **Micro-interactions** | Hover effects, transitions, loading indicators work |
| **Error states** | Submit invalid data → error message appears |
| **Console clean** | No JS errors, no 404s, no CSP violations |

### 3c. Capture Recording

Every browser verification creates a `.webp` recording artifact automatically.
Name recordings descriptively: `verify_job_listing`, `verify_bulk_upload`, etc.

**Gate**: All visual checks pass. No console errors.

---

## Layer 4: Final Server Restart

After every 3-step verification is complete, reset the application state to ensure clean operations going forward by running the restart server workflow:
@[.agent/workflows/restart-server.md]

---

## Completion Checklist

A feature is **DONE** when ALL boxes are checked:

- [ ] Unit test was written FIRST and **failed** (Red)
- [ ] Minimal code was written and test **passed** (Green)
- [ ] `go test -race ./...` passes (no regressions)
- [ ] `./scripts/pre-commit.sh` exits 0 (coverage ≥ 82.5%)
- [ ] `./scripts/verify_restart.sh` succeeds (server running)
- [ ] Browser subagent verified feature works for user
- [ ] Recording artifact saved with descriptive name
- [ ] `task.md` updated with completed status
- [ ] `spec.md` reviewed and updated if needed
- [ ] `@[.agent/workflows/restart-server.md]` was run after verification
- [ ] Commit with short, imperative message (e.g., "add bulk upload validation")
