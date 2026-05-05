# ADR: Runtime Observability via `verify trace`
**Date**: 2026-05-05 **Status**: Accepted

## 1. Context & User Problem
AI Agents (like Antigravity) often operate in a "blind spot" during the request/reponse lifecycle. When a handler fails or a database query returns unexpected results, the Agent must manually inject `fmt.Println` or `slog` statements, which is slow, prone to code pollution, and requires multiple edit-test cycles. This adds friction to the development loop and violates the goal of high-velocity feature delivery.

## 2. Decision
We are implementing a dedicated `verify trace` command and a corresponding `TraceMode` in the application server.
*   **TraceMode**: Triggered by `AGBALUMO_ENV=trace`, it injects `middleware.Logger()` and `middleware.BodyDump()` into the Echo server setup.
*   **`verify trace` CLI**: A wrapper that boots a temporary server in `TraceMode`, executes a targeted HTTP request (specified by `--path` and `--method`), captures the raw logs, and terminates the server.
*   **Shared Infrastructure**: Reuses the `BuildTestBinary` and `StartTestServer` logic from the Dynamic Audit suite to ensure environment parity.

## 3. The Complexity Kill-Switch (Rationale)
*   **User Value**: This is 10x better for Agentic development. Instead of guessing why an HTMX fragment failed to render, the Agent can run one command to see the exact input, output, and middleware logs.
*   **Performance Budget**: Zero impact on production. The tracing middleware is only active in `trace` mode.
*   **Minimalism Check**: This eliminates the need for "scratch logs" and temporary debugging prints that often leak into PRs.

## 4. Consequences
*   **Technical Tradeoffs**: Adds slight complexity to `server.go` and `config.go`, but this is isolated to the "Observability" domain.
*   **Observability**: This *is* the observability feature. It allows monitoring the request lifecycle on demand.
*   **SQLite Impact**: None, as it uses the existing test database configuration during the trace.

## 5. Alternatives Considered
**Manual Log Injection**: Rejected because it's slow, repetitive, and introduces risk of leaving debug code in production.
