# Echo Handler Patterns

## Overview
Standard patterns for Echo web handlers with HTMX integration, standardized error responses, and identity extraction.

## Handler Structure

### Basic Pattern
```go
func (h *XxxHandler) HandleEndpoint(c echo.Context) error {
    // 1. Extract identity (if needed)
    u, ok := user.GetUser(c)
    
    // 2. Extract parameters/Bind
    var req FormRequest
    if err := c.Bind(&req); err != nil {
        return ui.RespondError(c, err)
    }
    
    // 3. Call business logic (Service/Store)
    result, err := h.Service.DoWork(c.Request().Context(), u, req)
    
    // 4. Handle errors using UI helpers
    if err != nil {
        return ui.RespondError(c, err)
    }
    
    // 5. Return response (HTML or Fragment)
    return c.Render(http.StatusOK, "template", result)
}
```

### HTMX vs Full Page
```go
// Check for HTMX request
if c.Request().Header.Get("HX-Request") == "true" {
    // Return partial HTML fragment
    return c.Render(http.StatusOK, "partial_template", data)
}
// Return full page
return c.Render(http.StatusOK, "full_template", data)
```

## Error Handling

### Standardized Helpers
Always use `internal/ui` helpers for consistent UX:
*   `ui.RespondError(c, err)`: Renders an error page or toast.
*   `ui.RespondJSONError(c, code, message)`: Returns a structured JSON error.

## Identity Extraction
Always use `internal/module/user` helpers:
*   `user.GetUser(c)`: Safe extraction (returns `(*domain.User, bool)`).
*   `user.MustUser(c)`: Force extraction (panics if not found - use when middleware guarantees auth).

## Testing Patterns

### Integration Tests
Use `internal/testutil` for:
*   `testutil.SetupTestRepository(t)`: SQLite in-memory test DB.
*   `testutil.NewMainTemplate()`: Minimal templates for fragment testing.
*   `testutil.NewRealTemplate(t)`: Full filesystem templates for regression testing.
