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

3. **Flaky API Bad Gateway Mitigation**: If a `gh run` command fails with a transient `HTTP 502/503/504 Bad Gateway` server error from the GitHub API, do NOT abort the task. Wait 10 seconds and retry the query/command up to 3 times before raising a failure.

4. **Autonomous Hang Detection & Recovery (MANDATORY)**: E2E Playwright shards typically complete in ~1.5 to 2.5 minutes. If a remote Playwright job has been running or stuck in `Install Playwright Browsers & Deps` for **>10 minutes**, it is hung due to external download or runner stalling. Autonomously recover:
   - Cancel the hung run: `gh run cancel <run-id>`
   - Wait for cancellation status to stabilize: `gh run view <run-id>`
   - Rerun only the failed/cancelled shards: `gh run rerun <run-id> --failed`
   - Re-launch the background monitor: `gh run watch <run-id>`
   - Schedule a new 300s fail-safe timer, then yield.

5. If the run fails:
   - Identify the failed job and step.
   - Run `gh run view <run-id> --log-failed` to extract the traceback.
   - Fix and re-push. Repeat from Step 1.
   - Do NOT mark the task as complete until the remote CI passes.

## CI/CD Resiliency Standards

1. **Native Host Execution for E2E Tests**: To avoid transient container initialization overhead ("Initialize containers failed") and fragile nested dependency configurations (e.g., `setup-go` failing within a Docker container), configure Playwright E2E jobs to run natively on the runner host (`runs-on: ubuntu-latest`) rather than within a container block.
2. **Playwright Browser Caching**: To eliminate heavy downloads, always cache Playwright browser binaries in `~/.cache/ms-playwright` mapped to `package-lock.json` hashes. 
   - On cache hit: Invoke only `npx playwright install-deps` (installs OS system packages fast on the runner, ~15s).
   - On cache miss: Invoke `npx playwright install --with-deps`.
3. **Hard Step-Level Timeouts**: Always configure a strict `timeout-minutes: 5` on Playwright browser/dep setup steps to ensure immediate fail-fast behavior rather than 6-hour runner stalls.

## Post-Execution Health Check

1. Verify that the local server is running and healthy:
   `go run ./cmd/verify uptime`
   _Insight: This ensures the local dev environment is not left in a broken or stopped state after agent activity._

2. If it fails, restore the server using `go run ./cmd/verify watch`.
