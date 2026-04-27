---
name: CI Parity and Push Protocol
description: Ensure local CI parity with production and monitor remote CI execution.
triggers:
  - "pushing changes"
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

## Push & Remote Monitoring
1. Push changes to the remote branch.
2. Immediately watch the remote CI run using the GitHub CLI:
   `gh run watch`
3. Wait for the run to complete. If it fails:
   - Identify the failed job and step.
   - Run `gh run view <run-id> --log-failed` to extract the traceback.
   - Do NOT mark the task as complete until the remote CI passes.
