# Phase 10: Shadow Harness Purge

## Objective
Safely delete the remaining ad-hoc developer alias scripts (`scripts/*.sh`) that were created for the legacy "shadow harness" and Agentic orchestration but are no longer referenced or used by the codebase natively.

## Context
A global codebase scan reveals that several scripts have absolutely zero invocations anywhere in the codebase (neither in Go code, docs, nor GitHub Actions). These scripts include tools that Agents used to run manually (like `record-decision.sh` or `refactor.sh`). Because we rely entirely on the native Go workflow now with Phase 1-9 complete, these files are pure bloat.

## Steps for Execution
1. Open terminal and run the following deletion commands to permanently strip the shadow harness remnants:
   ```bash
   rm scripts/record-decision.sh
   rm scripts/refactor.sh
   rm scripts/ci-remote.sh
   rm scripts/repro_ci_failure.sh
   rm scripts/sandbox-setup.sh
   rm scripts/run_loadtest.sh
   ```
2. Optional Check: `scripts/benchmark_stress.sh` and `scripts/deploy_secrets.sh` are also completely unreferenced. Ask the user if they are run manually by human developers. If not, delete them as well.
3. Commit cleanly natively: `refactor(infra): purge unreferenced shadow harness scripts`.

## Verification
- The `scripts/` directory is significantly smaller and ONLY contains active hooks/entrypoints.
