# Handler Test Suite (internal/handler)

This package contains integration, latency, and benchmark test suites for HTTP routing.

## Note on Architecture
Application HTTP handlers and routes are organized into vertical-slice modules under `internal/module/` (`listing`, `admin`, `auth`, `feedback`, `user`). Refer to `internal/module/AGENTS.md` for handler implementation guidelines.

## Constraints
- **Latency Benchmark**: The search/filter endpoint must execute within budget (under 200ms without `-race`).
- **Isolation**: Use `cmd.SetupServer()` with test database configs to ensure realistic end-to-end execution.
