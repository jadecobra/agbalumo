# Development Guide

This document serves as the central source of truth for build, test, and quality control processes for Agbalumo.

## 🏗 Build & Run

### Safe Server Restart
The primary way to build and start the server is using the verification script:
```bash
./scripts/verify_restart.sh
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
A minimum test coverage threshold is required, enforced from `.agents/coverage-threshold`. This is checked by `scripts/pre-commit.sh`.

---

## 💎 Quality Control

### Pre-commit Hooks
Before committing any code, run the pre-commit script to ensure all gates pass:
```bash
./scripts/pre-commit.sh
```
This script performs:
1. Linting (`golangci-lint`)
2. Security scans (`gosec`)
3. Test execution and coverage checks
4. Performance audits

### Local CI Execution
To replicate the CI environment locally:
```bash
./scripts/ci-local.sh
```

### Performance Audit
We maintain a "Performance-First" culture. Run the automated audit to check asset sizes, DB config, and N+1 patterns:
```bash
./scripts/performance-audit.sh
```

---

## 🤖 Agentic Harness

Agbalumo uses an active operational framework for agentic coding.

### Workflow State Machine
Manage workflow states using `agent-exec.sh`:
```bash
./scripts/agent-exec.sh workflow status
```

### Workflow Gates
Programmatically verify workflow gates (e.g., Red tests) using:
```bash
./scripts/agent-gate.sh <gate_id>
```

### Brand Enforcement ("Juice")
Generate CSS tokens and Go constants from `.agents/rules/brand.toon`:
```bash
./scripts/generate-juice.sh
```

---

## 🔒 Security

Always run a security audit before any major release or dependency change:
```bash
./scripts/security-check.sh
```
