# Phase 8: Final Eradication of Legacy Bash Scripts

## Objective
Clean up the remaining fragmented shell scripts that are now broken or obsolete due to the removal of `Taskfile.yml` and the consolidation into `cmd/verify`.

## Context
Various scripts still rely on `task` or have been entirely subsumed by the Pure Go architecture. This completes the transition.

## Steps for Execution
1. Update `scripts/watch.sh`: It currently executes `exec task watch`. Change it to `exec go run cmd/verify/main.go watch` (if implemented) OR configure `air` or `reflex` properly. Or simply delete it if watching is no longer needed via script.
2. Update `scripts/setup-hooks.sh`: Replace its injection of `task pre-commit` with `go run cmd/verify/main.go precommit` in the `.git/hooks/pre-commit` file generated.
3. Update `scripts/verify_restart.sh`: Change `task pre-commit` and `task restart` to direct commands (`go run cmd/verify/main.go precommit`).
4. Update `scripts/browser_audit.sh`: Replace `task restart-server`.
5. Audit and Delete obsolete files:
   - `scripts/gitleaks-scan.sh` (Obsolete, natively handled)
   - `scripts/utils.sh` (If Phase 7 removes the last dependency)
   - `scripts/utils/audit_helpers.sh` (If Phase 6 removes the dependency)
6. Commit all changes natively: `refactor(infra): eradicate final taskfile dependencies and broken bash scripts`.

## Verification
- Running `git grep "task "` inside `scripts/` should yield zero pipeline invocations.
- All deprecated bash utilities have been successfully deleted.
