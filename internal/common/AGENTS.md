# Common Package (internal/common)

Shared utilities for error toasts and static page handlers.

## Constraints
- Canonical error rendering is `ui.RespondError(c, err)` from `internal/ui`.
- `RenderImageErrorToast` renders HTMX out-of-band error toasts for upload failures.
- `PageHandler` serves static pages (e.g., `/about`). It uses `BaseHandler`-style rendering from `internal/module/base.go`.
