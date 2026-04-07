# Task 15: CLI Command Abstraction (Verify Tool)

## Context
The `cmd/verify/main.go` file contains many duplicated `exec.Command` calls with identical error handling logic. This task creates shared internal helpers to handle command execution and results.

## Checklist
- [ ] Implement `runCmd(name string, args ...string) error` for shared execution logic.
- [ ] Refactor `api-spec`, `template-drift`, and `precommit` commands to use this helper.
- [ ] Abstract common flag parsing for `race` and `threshold-path`.

## Verification
- [ ] Run `go run cmd/verify/main.go ci` and ensure the tool is still functional.
- [ ] Verify that duplicate code was significantly reduced in `cmd/verify/main.go`.
