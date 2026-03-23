# Build/Lint/Test Commands

### Run All Tests
```bash
go test -json ./...
```

### Run Single Test
```bash
go test -json -v -run TestFunctionName ./internal/package/
go test -json -v -run TestFunctionName/SubtestName ./internal/package/
```

### Run Tests with Race Detection
```bash
go test -json -race ./...
```

### Pre-Commit Quality Gate
```bash
./scripts/pre-commit.sh
```
Runs: `gofmt`, `go mod tidy`, `go vet`, race tests, coverage (threshold in `.agent/coverage-threshold`), secret scanning.

### Build & Restart Server
```bash
./scripts/verify_restart.sh
```
Builds CSS, compiles binary, restarts server on :8443.

### Build CSS
```bash
npm run build:css
```
