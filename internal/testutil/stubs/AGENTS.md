# Test Stubs Context

Shared mock implementations and test doubles.

## Usage Standards
*   **Simplicity**: Stubs should be minimal. If a mock requires complex state management, use a real in-memory SQLite instance instead.
*   **Package Boundary**: Do not use these stubs in production code. They are strictly for `*_test.go` files.
*   **Consistency**: Ensure stubs implement the interfaces defined in `internal/domain/`.

## Common Stubs
* `StubRepository`: In-memory map-based store.
* `StubGeocoder`: Returns fixed coordinates for testing.
