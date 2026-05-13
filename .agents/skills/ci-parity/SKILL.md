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

1. **Determine Execution Path:**
   - **Fast Path (Diff-based/Sharded):** If you only changed specific UI components or templates, run Playwright strictly for those impacted routes using `--focus` or equivalent test filters: `go run ./cmd/verify ci --with-docker --focus=visual.spec.ts`.
   - **Full Path (Architecture/Core Changes):** If you modified core `internal/domain` logic, dependencies, or Dockerfiles, run the full suite: `go run ./cmd/verify ci --with-docker`.

2. Fix any local violations before pushing.
   *Insight: Running in Docker ensures we catch environment drift and security vulnerabilities before production, but we must use target filters to preserve velocity.*

## Push & Remote Monitoring
1. Execute the push and automated monitoring wrapper:
   `./scripts/pushw.sh`
   *Insight: This atomically executes the push and polls the GitHub API for the specific commit's CI run ID to resolve race conditions.*

2. Manual Fallback (if the script fails or is bypassed):
   `gh run watch`

3. If the run fails:
   - Identify the failed job and step.
   - Run `gh run view <run-id> --log-failed` to extract the traceback.
   - Do NOT mark the task as complete until the remote CI passes.

## Post-Execution Health Check
1. Verify that the local server is running and healthy:
   `go run ./cmd/verify uptime`
   *Insight: This ensures the local dev environment is not left in a broken or stopped state after agent activity.*

2. If it fails, restore the server using `go run ./cmd/verify watch`.
