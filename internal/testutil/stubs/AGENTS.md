# Test Stubs Context

Shared mock implementations and test doubles.

## Usage Standards
*   **Simplicity**: Stubs should be minimal. If a mock requires complex state management, use a real in-memory SQLite instance instead.
*   **Package Boundary**: Do not use these stubs in production code. They are strictly for `*_test.go` files.
*   **Consistency**: Ensure stubs implement the interfaces defined in `internal/domain/`.

## Common Stubs
* `ErrDefault`: Standard test stub error sentinel.
* `NoopLogger`: Discards all log statements.
* `FailingLogger`: Panics on log invocation for assertion paths.
* `StubResolver`: Generic resolver double accepting a custom `ResolveFunc`.
