# Technical Specification: Taskfile vulncheck Optimization

## Objective
Optimize the `vulncheck` and `ci:vulncheck` tasks to use a pre-installed binary if available, avoiding redundant `go install` or `go run` overhead.

## Proposed Changes

### Taskfile.yml Architecture

#### New Task: `vulncheck:install`
- **Purpose**: Idempotently install the `govulncheck` binary into the workspace-local `GOPATH`.
- **Location**: `.tester/tmp/go/bin/govulncheck`
- **Internal**: `true`
- **Prerequisites**: `status` check using `test -f ./.tester/tmp/go/bin/govulncheck`.

#### Modified Task: `vulncheck`
- **Dependency**: Add `vulncheck:install` to `deps`.
- **Command**: Run `{{.TASKFILE_DIR}}/.tester/tmp/go/bin/govulncheck ./...` directly.

#### Modified Task: `ci:vulncheck`
- **Dependency**: Add `vulncheck:install` to `deps`.
- **Command**: Replace `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` with the direct binary call.

## API Contracts / CLI
No external contract changes. This is a performance and logic optimization.

## Security Considerations
The `govulncheck` tool remains the source of truth for vulnerability scanning. Using a local binary avoids network access on subsequent runs, reducing CI/CD surface area and potential flakiness.

## Verification
1. Remove `.tester/tmp/go/bin/govulncheck` then run `task vulncheck`. Verify installation occurs.
2. Run `task vulncheck` again. Verify `go install` is skipped.
3. Run `task ci:vulncheck`. Verify direct binary usage.
