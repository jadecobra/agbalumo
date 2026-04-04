# Phase 4: Final Infrastructure Decommission

## Objective
Safely delete the outdated `Taskfile.yml` and obsolete CI bash scripts now that `cmd/verify/main.go` is the absolute source of truth.

## Context
With dev/prod parity secured via the Go pipeline, we eradicate tech debt and eliminate the possibility of a developer running legacy, drifting bash execution paths.

## Steps for Execution
1. Update `.git/hooks/pre-commit` (or the script defining it, e.g. `scripts/setup-hooks.sh`). It should explicitly replace `task pre-commit` with `go run cmd/verify/main.go precommit`.
2. Update the system's primary GitHub Actions CI yaml file (e.g. `.github/workflows/ci.yml`) to replace any `task ci` or `task ci:*` commands with `go run cmd/verify/main.go ci`.
3. Delete the `Taskfile.yml` entirely (`rm Taskfile.yml`).
4. Audit `scripts/` directory and delete all scripts that were directly replicated into Go commands, including but not limited to:
   - `scripts/verify-ci-tools.sh`
   - `scripts/critique.sh` (if the critique command was moved or is deprecated)
   - `scripts/verify-golangci-config.sh`
5. Ensure `cmd/verify/main.go` builds natively: `go build -o .bin/verify cmd/verify/main.go`.
6. Commit all deletions. Message: `refactor(infra): decommission Taskfile and bash ci orchestrations`.

## Verification
- The pipeline CI in GitHub Actions should remain perfectly green.
- Local commits should trigger the native Go hook without any dependency on the legacy `Taskfile.yml`.
