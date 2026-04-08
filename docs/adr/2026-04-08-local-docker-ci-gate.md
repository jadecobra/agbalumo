# ADR: Local Docker Build Gate in CI

**Date**: 2026-04-08 **Status**: Accepted

## Context
Broken Dockerfile stages reached production multiple times because local CI lacked a container build step. The gap between local Go-native CI and production Docker-based CI allowed regressions in the build process (e.g., Litestream builder failures) to bypass local validation.

## Decision
Add an opt-in `--with-docker` flag to `go run cmd/verify/main.go ci` that appends a `docker build` step as personal validation. While not mandatory for all runs to preserve performance, it is established as a mandatory gate for any changes touching the `Dockerfile`.

## Consequences
- **Correctness**: Dockerfile regressions are caught before they reach production CI.
- **Performance**: Validating with Docker adds 3-5 minutes to the CI run; hence the opt-in flag.
- **Dependencies**: Developers must have Docker installed locally to validate Docker-related changes.
- **Gap Reduction**: Local CI now better mirrors production CI constraints.