# Internal/Handler - HTTP Handler Patterns

## Overview
Echo web handlers with HTMX integration, SQLite persistence, and image processing.

## Handler Structure

### Basic Pattern
```go
func (h *XxxHandler) HandleEndpoint(c echo.Context) error {
    // 1. Extract parameters
    // 2. Validate input
    // 3. Call service layer
    // 4. Handle errors
    // 5. Return response
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

## Validation Patterns

### Echo Validation
```go
// Use echo.Context.Validate() for struct validation
err := c.Validate(inputStruct)
if err != nil {
    return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Validation Error: "+err.Error()))
}
```

### Custom Validation
```go
// Domain-level validation
listing, err := ToListing(input)
if err != nil {
    return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, err.Error()))
}
```

## Error Handling

### Centralized Response
```go
// Use RespondError wrapper for all handlers
return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Resource not found"))
```

### Error Types
- **400 Bad Request**: Validation errors, missing parameters
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Unexpected errors (logged internally)

## File Upload Patterns

### Image Handling
```go
// Handle file upload with validation
file, err := c.FormFile("image")
if err != nil {
    return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Missing image file"))
}

// Process and save
imageURL, err := h.ImageService.UploadImage(file)
if err != nil {
    return RespondError(c, err)
}
```

## Authentication

### User Context
```go
// Get authenticated user from middleware context
user, ok := c.Get("User").(*domain.User)
if !ok || user == nil {
    return RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required"))
}
```

## Response Formatting

### JSON Responses
```go
// For API endpoints
return c.JSON(http.StatusOK, data)
```

### HTML Responses
```go
// For web pages
return c.Render(http.StatusOK, "template_name", data)
```

## Testing Patterns

### Table-Driven Tests
```go
func TestHandler(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        expectErr bool
    }{
        {name: "valid case", input: "valid", expectErr: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Mock Usage
```go
// Use testify/mock for service layer
mockService := new(mock.Service)
mockService.On("FindByID", ctx, "123").Return(domain.Listing{}, nil)
```

## Common Patterns

### URL Normalization
```go
// Normalize URLs for consistency
func normalizeURL(url string) string {
    // ...
}
```

### File Headers
```go
// Get file header for uploads
header, err := getFileHeader(file)
if err != nil {
    return RespondError(c, err)
}
```

## Anti-Patterns

- **No direct database calls** - Always use service layer
- **No raw SQL** - Use repository patterns
- **No error suppression** - Always return or log errors
- **No direct template rendering** - Use renderer service
- **No missing validation** - Always validate input

## Dependencies
- **echo/v4** - Web framework
- **testify/mock** - Test mocking
- **image service** - File upload processing
- **auth middleware** - User context
- **renderer service** - Template rendering

## Coverage Requirements
- All handler functions must have tests
- Error paths must be tested
- HTMX vs full page logic must be tested
- File upload edge cases must be tested

## Testing Commands
```bash
go test -json -v -run TestHandlerName ./internal/handler/
go test -json -v -run TestHandlerName/SubtestName ./internal/handler/
```

## Notes
- Always restart server after handler changes
- Use RespondError for consistent error responses
- Test both success and error paths thoroughly
