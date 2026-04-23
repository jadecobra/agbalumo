# HTTP Handlers: Agent Guidance

This package manages the HTTP layer using the Echo framework.

# Handler Constraints  
- Use `RespondError(c, err)` — never raw `c.JSON()`
- All form bindings use `form` struct tags
- No raw HTML in handlers — use `ui/templates/components/`

## Handling Principles

- **Thin Handlers**: Handlers should only bind input, call services, and render output. Business logic belongs in `internal/service/`.
- **Fail Fast**: Validate all required form fields or parameters using `c.FormValue` or binding structs before calling downstream services.
- **Error Mapping**: Map all service errors to appropriate HTTP status codes (e.g., `domain.ErrNotFound` -> `http.StatusNotFound`).

## UI & Rendering

- **HTMX Conscious**: Many handlers return HTMX fragments or trigger OOB (Out-of-Band) swaps. Check for the `HX-Request` header.
- **Template Semantic Tags**: Always add `data-template-file="[relative/path/to/template.html]"` to the root element of any rendered template to aid debugging.

## Common Pitfalls

- **Context Handling**: Always pass `c.Request().Context()` to downstream services, never `context.Background()`.
- **Session Management**: Use the centralized session keys defined in `internal/domain/constants.go`.
