# Development Guide

This document serves as the central source of truth for build, test, and quality control processes for Agbalumo.

## 🏗 Build & Run

### Safe Server Restart
The primary way to build and start the server is using the verification script:
```bash
go build
```
This script handles process cleanup, environment variable validation, and secure HTTPS setup.

### HTTPS & Certificates
The development server runs on **Port 8443**. Self-signed certificates are located in `certs/`.

---

## 🧪 Testing Strategy

We follow a strict **TDD (Red-Green-Refactor)** protocol.

### Running Tests
```bash
# Run all tests with race detection
go test -json -v -race ./...

# Run tests with coverage
mkdir -p .tester/coverage
go test -json -coverprofile=.tester/coverage/coverage.out ./...
go tool cover -func=.tester/coverage/coverage.out

# Generate HTML coverage report
go tool cover -html=.tester/coverage/coverage.out -o .tester/coverage/coverage.html
```

### Coverage Threshold
A minimum test coverage threshold is required, enforced from `.agents/coverage-threshold`. This is checked by `task pre-commit`.

---

## 💎 Quality Control

### Pre-commit Hooks
Before committing any code, run the pre-commit script to ensure all gates pass:
```bash
task pre-commit
```
This script performs:
1. Linting (`golangci-lint`)
2. Security scans (`gosec`)
3. Test execution and coverage checks
4. Performance audits

### Local CI Execution
To replicate the CI environment locally:
```bash
go run ./cmd/verify ci
```

### Performance Audit
We maintain a "Performance-First" culture. Run the automated audit to check asset sizes, DB config, and N+1 patterns:
```bash
go run ./cmd/verify perf
```

---



---

## 🔒 Security

Always run a security audit before any major release or dependency change:
```bash
go run ./cmd/verify audit
```
