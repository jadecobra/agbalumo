# Task 24c: Reduce Complexity — `(*AdminHandler).HandleAddCategory` (score: 13 → target ≤ 10)

## File
`internal/module/admin/admin.go`

## Context

`HandleAddCategory` scores **13** due to deeply nested logic inside the `if err == nil`
block when checking for duplicate categories, as well as repeating the boilerplate to
add flash messages to the session.

The logic blocks driving the score:
1. `if name == ""` (early exit)
2. `if err == nil` (checking if fetching categories succeeded)
3. `for _, cat := range existing` (iterating categories)
4. `if strings.EqualFold(cat.Name, name)` (checking match)
5. `if sess != nil` (flash message boilerplate repeated twice in function)

## Current Code (lines 232–283)

```go
func (h *AdminHandler) HandleAddCategory(c echo.Context) error {
    // ... setup ...
    name := strings.TrimSpace(c.FormValue("name"))
    if name == "" {
        return c.Redirect(http.StatusFound, "/admin")
    }

    existing, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
    if err == nil {
        for _, cat := range existing {
            if strings.EqualFold(cat.Name, name) {
                sess := customMiddleware.GetSession(c)
                if sess != nil {
                    sess.AddFlash(fmt.Sprintf("Category '%s' already exists!", cat.Name), "message")
                    _ = sess.Save(c.Request(), c.Response())
                }
                return c.Redirect(http.StatusFound, "/admin")
            }
        }
    }
    // ... creates category, then repeats the flash message block for success ...
}
```

## Required Changes

**Step 1**: Add a simple helper to `internal/module/admin/admin.go` for finding duplicates:

```go
func hasDuplicateCategory(existing []domain.CategoryData, name string) bool {
    for _, cat := range existing {
        if strings.EqualFold(cat.Name, name) {
            return true
        }
    }
    return false
}
```

**Step 2**: Add a helper for setting flash messages and redirecting:

```go
func flashAndRedirect(c echo.Context, msg, url string) error {
    if sess := customMiddleware.GetSession(c); sess != nil {
        sess.AddFlash(msg, "message")
        _ = sess.Save(c.Request(), c.Response())
    }
    return c.Redirect(http.StatusFound, url)
}
```

**Step 3**: Simplify `HandleAddCategory` leveraging those helpers:

```go
func (h *AdminHandler) HandleAddCategory(c echo.Context) error {
    ctx := c.Request().Context()
    name := strings.TrimSpace(c.FormValue("name"))
    if name == "" {
        return c.Redirect(http.StatusFound, "/admin")
    }

    if existing, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false}); err == nil {
        if hasDuplicateCategory(existing, name) {
            return flashAndRedirect(c, fmt.Sprintf("Category '%s' already exists!", name), "/admin")
        }
    }

    claimable := c.FormValue("claimable") == "true"
    now := time.Now()
    cat := domain.CategoryData{
        ID:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
        Name:      name,
        Claimable: claimable,
        Active:    true,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := h.App.DB.SaveCategory(ctx, cat); err != nil {
        c.Logger().Errorf("failed to save custom category: %v", err)
    }

    return flashAndRedirect(c, "Category added successfully!", "/admin")
}
```

## What NOT to change
- Do not modify other handlers in `admin.go`
- Do not export the new helpers

## Verification

```bash
go test ./internal/module/admin/...
go run cmd/verify/main.go critique 2>&1 | grep "HandleAddCategory"
```

The function must no longer appear in the complexity output, or score ≤ 10.

## Commit
```
refactor(admin): extract helpers to reduce HandleAddCategory complexity
```
