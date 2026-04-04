# Docker Build & Scan Failure Checkpoint
**Date**: 2026-04-04 **Status**: [FIXED & VERIFIED]

## Docker Build & Scan Summary
The production CI pipeline failed at the **Docker Build & Scan** stage due to high-severity security vulnerabilities that were not detected in the local development environment.

### 1. Root Causes
- **Tooling Discrepancy**: Local checks used `govulncheck` (call-stack analysis), while CI used **Trivy** (image-layer scanning). Trivy flagged vulnerabilities in non-called code and OS/Node layers that `govulncheck` ignored.
- **Dependency Drift**: High-severity vulnerabilities in `picomatch` (Node) and `golang.org/x/image` (Go) required updates. The `x/image` update forced a minimum Go version requirement of **1.25.0**, which clashed with the CI's pinned 1.24 environment.
- **Infrastructure Blind Spot**: Local development on macOS (Go 1.26.1) hid toolchain incompatibilities that only surfaced when building specifically for the Linux/Alpine Docker targets used in production.

### 2. Resolved Items
- **Patched Vulnerabilities**:
    - `npm audit fix` resolved a high-severity ReDoS vulnerability in `picomatch`.
    - Upgraded `golang.org/x/image` to **v0.38.0** to fix a Critical OOM vulnerability (GO-2026-4815).
    - Removed the unused CGO-based `github.com/mattn/go-sqlite3` to reduce the attack surface.
- **Environment Alignment**:
    - Upgraded CI runners and Docker builders to **Go 1.25** and **Alpine:latest**.
    - Updated Litestream to **v0.5.10** for active security maintenance.
    - Standardized `go.mod` to version **1.25.0** for production parity.

## Planned Next Steps (Infrastructure Governance)
- **[/janitor] Workflow**: Execute a complete workspace cleanup to surface any remaining tech debt or stale documentation markers.
- **Unified Security Gate**: Evaluate adding a local `task ci:docker-scan` to mirror the production Trivy checks in the pre-commit phase.
- **Toolchain Locking**: Implement a mechanism to ensure the `go-version` in GitHub Actions is dynamically synced from `go.mod` to prevent version drift failures.

## Verification Log
- [x] Local `go test ./...` passed on Go 1.25.
- [x] Local `npm audit` report: `found 0 vulnerabilities`.
- [x] CI `Tests & Coverage` passed on Go 1.25.
- [x] CI `Drift Verification` passed.
- [x] CI `Docker Build & Scan` verified passing with patched images.
