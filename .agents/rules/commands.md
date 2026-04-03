# Build/Lint/Test Commands

### Squad Harness (`agent-exec.sh`)
The central interface for feature development and phase management.
```bash
./scripts/agent-exec.sh init <feature> <workflow>  # Start new feature
./scripts/agent-exec.sh set-phase <phase>           # Transition phase (RED/GREEN/REFACTOR)
./scripts/agent-exec.sh verify <gate_id>             # Pass a specific gate (test/lint/security)
./scripts/agent-exec.sh handoff <persona>            # Execute persona handoff
./scripts/agent-exec.sh cost                         # Report context token cost
```

### Verification Scripts
```bash
go run scripts/verify-persona.go  # Validate persona configurations
./scripts/ci-local.sh             # Run full CI suite locally
```

### Run All Tests
```bash
go test -json ./...
```

### Run Single Test
```bash
go test -json -v -run TestFunctionName ./internal/package/
```

### Pre-Commit Quality Gate
```bash
task pre-commit
```
Runs: `gofmt`, `go mod tidy`, `go vet`, race tests, token-based cost check, secret scanning.

### Build & Restart Server
```bash
./scripts/verify_restart.sh
```
Builds CSS, compiles binary, restarts server on :8443.
