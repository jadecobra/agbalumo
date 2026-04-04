# Phase 7: Gate Enforcement Migration

## Objective
Convert `scripts/test_gate_enforcement.sh` and its underlying functions in `scripts/utils.sh` to a native Go `check-gates` subcommand inside `cmd/verify/main.go`. 

## Context
The project uses strict TDD/Agent pipeline gates (RED/GREEN/REFACTOR phases stored in local JSON). The shell scripts assert that these states are valid. Moving this logic to Go eliminates JSON parsing issues with `jq` and brittle environment boundaries.

## Steps for Execution
1. Open `cmd/verify/main.go`.
2. Add a `checkGatesCmd` Cobra command.
3. In `checkGatesCmd`, implement logic to:
   - Read the agent state JSON (likely located in an `.agents/state.json` or `.metrics/` directory, depending on where `utils.sh` looks).
   - Use Go's `encoding/json` to Unmarshal the state.
   - Enforce the rules natively (e.g., if Phase == "RED", ensure `gates["red-test"] == "PASS"`).
4. Register the new command.
5. Port the tests from `scripts/test_gate_enforcement.sh` into `cmd/verify/main_test.go` or a similar dedicated Go test file to maintain coverage.
6. Delete `scripts/test_gate_enforcement.sh`.
7. Commit changes natively: `refactor(ci): migrate workflow gate enforcement to pure go`.

## Verification
- `scripts/test_gate_enforcement.sh` is completely deleted, and `go run cmd/verify/main.go check-gates` operates intelligently using native JSON structs.
