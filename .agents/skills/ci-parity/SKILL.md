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
1. Run full local CI suite mirroring production:
   `go run ./cmd/verify ci --with-docker`
   *Insight: This catches environment drift, Docker build failures, and Trivy security vulnerabilities before they reach production.*

2. Fix any local violations before pushing.
2b. If any UI/template/Go files changed, run `go run ./cmd/verify ci --with-docker` — this runs Playwright inside a Linux Docker container to regenerate and validate linux snapshots before push.

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
