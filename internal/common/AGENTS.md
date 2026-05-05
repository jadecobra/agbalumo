# Common Package (internal/common)

Shared utilities for error rendering and static page handlers.

## Constraints
- `RespondError(c, err)` is the canonical error rendering function. All modules must use it.
- `PageHandler` serves static pages (e.g., `/about`). It uses `BaseHandler`-style rendering from `internal/module/base.go`.
