# Build Tools Context

Handles build-time dependencies and code generation.

## Constraints
*   **Isolation**: This package uses the `// +build tools` pattern. It MUST NOT be imported by any runtime code.
*   **Version Pinning**: All tools listed here are pinned in `go.mod` to ensure reproducible builds.

## Purpose
* Pins development CLI tools: `golangci-lint`, `goconst`, `dupl`, `gocognit`, `gitleaks`, `fieldalignment`, and `govulncheck`.
