# Technical Specification: Taskfile Tooling Optimization

## Objective
Extend the performance optimization of `vulncheck` (binary existence checks) to `golangci-lint` and `gitleaks` to reduce local and CI/CD overhead and ensure consistent architecture-agnostic installation.

## Proposed Changes

### Taskfile.yml Architecture

#### New Task: `lint:install`
- **Purpose**: Idempotently install `golangci-lint` (`v1.64.5`) to `.tester/tmp/go/bin/golangci-lint`.
- **Internal**: `true`
- **Prerequisites**: `status` check using `test -f ./.tester/tmp/go/bin/golangci-lint`.
- **Install Command**: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5`.

#### New Task: `gitleaks:install`
- **Purpose**: Idempotently install `gitleaks` (`v8.21.2`) to `.tester/tmp/go/bin/gitleaks`.
- **Internal**: `true`
- **Prerequisites**: `status` check using `test -f ./.tester/tmp/go/bin/gitleaks`.
- **Install Command**: `go install github.com/zricethezav/gitleaks/v8@v8.21.2`.

#### Modified Task: `lint` & `ci:lint`
- **Dependency**: Add `lint:install` to `deps`.
- **Command**: Run `{{.TASKFILE_DIR}}/.tester/tmp/go/bin/golangci-lint run ...` directly.

#### Modified Task: `gitleaks`
- **Dependency**: Add `gitleaks:install` to `deps`.
- **Command**: Update call to `.tester/tmp/go/bin/gitleaks`.

### script/gitleaks-scan.sh
- **Refactor**: Prioritize local binary path during `setup_path` or explicit binary checks.
- **Simplification**: Remove slow `curl`/`brew` logic if the `gitleaks:install` task is guaranteed to have run (via `Taskfile` dependency).

## API Contracts / CLI
No external contract changes.

## Security Considerations
Pinned versions for `golangci-lint` and `gitleaks` ensure deterministic security scans across environments and prevent supply-chain drift.

## Verification
1. Run `task lint`. Monitor "installing" behavior only on first run.
2. Run `task gitleaks`. Verify local binary usage.
3. Update `internal/agent/task_optimization_test.go` (renamed from `task_vulncheck_test.go`) to verify all three tools.
