---
name: CI Parity and Push Protocol
description: Ensure local CI parity with production and monitor remote CI execution.
triggers:
  - "push changes"
  - "CI failure"
  - "production parity"
mutating: false
---

# CI Parity & Push Protocol Skill

## Local Verification (Pre-Push)

To avoid CI Matrix Bloat and exorbitant local time costs, you MUST NOT run the entire test suite blindly if only a subset of files changed.

1. **Scope the diff first (MANDATORY):**
   ```bash
   git diff --name-only HEAD~1 HEAD
   ```
   Use the output to select the execution path:
   - **Fast Path:** Only `ui/templates/`, `ui/static/`, or `internal/handler/` files changed → run targeted tests:
     `go run ./cmd/verify ci --with-docker --focus=visual.spec.ts`
   - **Full Path:** Any `internal/domain/`, `internal/repository/`, `go.mod`, `Dockerfile`, or `.github/` files changed → run the full suite:
     `go run ./cmd/verify ci --with-docker`
   - **Rule:** When in doubt, use Fast Path first. Escalate to Full Path only if Fast Path reveals a dependency failure.

2. Fix any local violations before pushing.

## Linux Snapshot Synchronization (MANDATORY after any UI change touching snapshots)

Visual snapshots are platform-specific. Darwin snapshots generated locally will always diverge from the Linux CI environment. Run this after any change that affects snapshot-covered templates (`sandbox.html`, `modal_detail.html`):

```bash
# ARM64 Mac (Apple Silicon) — matches CI runner architecture
# Resolve Playwright version dynamically from package.json to prevent drift
PW_VER=$(node -e "console.log(require('./package.json').devDependencies['@playwright/test'].replace(/[\^~]/,''))")
GOOS=linux GOARCH=arm64 go build -o server-linux main.go && \
docker run --rm -v $(pwd):/app -w /app \
  -e "AGBALUMO_TEST_SERVER_COMMAND=./server-linux serve" \
  -e "AGBALUMO_ENV=test" \
  "mcr.microsoft.com/playwright:v${PW_VER}-noble" \
  sh -c "npx playwright test visual.spec.ts --update-snapshots" && \
rm server-linux
```
_Insight: The version is sourced from `package.json` so the image tag automatically tracks Playwright version bumps without requiring a manual edit to this skill._

Then verify parity before committing:

```bash
go run ./cmd/verify snapshot-parity
```

## Push & Remote Monitoring

1. Execute the push and automated monitoring wrapper:
   `./scripts/pushw.sh`
   _Insight: This atomically executes the push and polls the GitHub API for the specific commit's CI run ID to resolve race conditions._

2. **Reactive Monitoring with Fail-Safe (MANDATORY):** After the push, do NOT poll using short-interval timers. Instead, launch the monitoring command in the background:
   ```bash
   gh run watch <run-id> --exit-status
   ```
   Set a single **300-second (5 minutes)** fail-safe timer via the `schedule` tool to capture hangs, then immediately yield the turn. The system's high-priority reactive completion will automatically wake you up with the final status when the pipeline concludes.
   - Exit code `0` = all jobs passed. You may proceed.
   - Any non-zero exit = pipeline failed. Do NOT declare the task complete; analyze failure logs via `gh run view <run-id> --log-failed` and resolve.

3. If the run fails:
   - Identify the failed job and step.
   - Run `gh run view <run-id> --log-failed` to extract the traceback.
   - Fix and re-push. Repeat from Step 1.
   - Do NOT mark the task as complete until the remote CI passes.

## CI/CD Resiliency Standards

1. **Native Host Execution for E2E Tests**: To avoid transient container initialization overhead ("Initialize containers failed") and fragile nested dependency configurations (e.g., `setup-go` failing within a Docker container), configure Playwright E2E jobs to run natively on the runner host (`runs-on: ubuntu-latest`) rather than within a container block.
2. **Dynamic Browser Provisioning**: Install browsers and system-level library dependencies natively on the host using `npx playwright install --with-deps`. This is faster and immune to Docker daemon startup or pulling delays.

## Post-Execution Health Check

1. Verify that the local server is running and healthy:
   `go run ./cmd/verify uptime`
   _Insight: This ensures the local dev environment is not left in a broken or stopped state after agent activity._

2. If it fails, restore the server using `go run ./cmd/verify watch`.
