# Task Checkpoint: CI Pipeline Disparity and Visual Regression

## Current State
The objective is to fix persistent Playwright UI test failures on GitHub Actions, specifically the `Listing Modal Visual Regression` tests for `Desktop` and `Wide` viewports. 

1. **Snapshots Updated:** I generated fresh `linux/amd64` baseline snapshots locally using Docker (`mcr.microsoft.com/playwright:v1.59.1-noble`) to match the CI environment's architecture. However, pushing these updates did not resolve the GitHub Actions test failure.
2. **CI Pipeline Hardened:** Discovered that the remote CI workflow was running Playwright directly on `ubuntu-latest` rather than inside the Docker container, leading to potential rendering variations. Updated `.github/workflows/ci.yml` to explicitly run the `e2e_tests` job inside the Playwright Docker container and configured it to upload the `playwright-report` artifacts upon failure.
3. **Artifact Retrieval:** The test failed again, but we successfully downloaded the generated `playwright-report` from GitHub Actions.
4. **Local Investigation:** The `playwright-report` zip was extracted locally, and an HTTP server was spun up on port `8080`. The user currently has the Playwright Report open in their local browser (`http://localhost:8080/`) to manually inspect the visual diffs.

## Errors Encountered
- **Visual Regression Flakes:** The `Listing Modal Visual Regression` test continues to fail on GitHub Actions despite syncing the `linux/amd64` baseline snapshots.
- **Missing CI Artifacts:** Previously, failing CI runs did not export the `playwright-report`, making it impossible to see *why* the visual tests were failing without running them locally. This has now been resolved.
- **Docker Architecture Mismatches:** Initial attempts to regenerate Linux snapshots on a Mac M-series machine defaulted to `linux/arm64`, which fails against the GitHub Actions `linux/amd64` environment. This was mitigated by using `--platform linux/amd64` during local snapshot generation.

## Planned Next Steps
1. **Analyze the Diffs:** The user (or the agent via browser subagent) needs to review the exact visual discrepancies shown in the `playwright-report` running at `http://localhost:8080/`. We need to identify if the failure is due to a shift in content, a rendering quirk (e.g., antialiasing, missing font), or a dynamic element loading asynchronously.
2. **Implement Fix:** Depending on the diff analysis:
   - If it's an asynchronous loading issue, increase the timeout or add wait assertions in `tests/e2e/visual.spec.ts` (e.g., `await page.waitForTimeout(1000)` or wait for a specific element).
   - If it's a CSS/rendering issue (e.g., scrollbars, transparent backgrounds), patch the UI elements to be deterministic or apply a `maxDiffPixelRatio` tolerance in Playwright.
3. **Verify Locally and Push:** Regenerate the snapshots locally using the `linux/amd64` docker container if needed, run `go run ./cmd/verify precommit`, and push to ensure the CI pipeline finally passes cleanly.
