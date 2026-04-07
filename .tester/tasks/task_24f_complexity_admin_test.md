# Task 24f: Reduce Complexity — `TestAdminHandler_HandleToggleFeatured` (score: 14 → target ≤ 10)

## File
`internal/module/admin/admin_actions_test.go`

## Context

`TestAdminHandler_HandleToggleFeatured` scores **14** because it is a table-driven test
with a heavy inline `for _, tt := range tests` loop. Inside the loop, it performs URLs
setup, authentication setup, routing execution, and multiple conditional assertions.

## Current Code (lines 97–131)

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        formData := url.Values{}
        formData.Set("featured", tt.featured)
        urlPath := "/admin/listings/" + tt.id + "/featured"
        if tt.id == "" {
            urlPath = "/admin/listings/featured"
        }
        c, rec := setupAdminTestContext(http.MethodPost, urlPath, strings.NewReader(formData.Encode()))
        setupAdminAuth(t, c)
        if tt.id != "" {
            c.SetParamNames("id")
            c.SetParamValues(tt.id)
        }

        app, h, cleanup := setupAdminTest(t)
        defer cleanup()
        tt.setupData(t, app.DB)

        _ = h.HandleToggleFeatured(c)
        assert.Equal(t, tt.expectCode, rec.Code)

        if tt.expectCode == http.StatusOK {
            htmlResponse := rec.Body.String()
            assert.Contains(t, htmlResponse, "listing-row-")
            assert.NotContains(t, htmlResponse, "{\"featured\":")
        }

        if tt.expectCode == http.StatusOK && tt.id == "123" {
            testutil.AssertFeaturedStatus(t, app.DB, tt.id, true)
        }
    })
}
```

## Required Changes

**Step 1**: Move the HTMX/HTML assertions into a helper function inside the file:

```go
func assertFeaturedResponse(t *testing.T, rec *httptest.ResponseRecorder, expectCode int, id string, db domain.ListingRepository) {
    t.Helper()
    assert.Equal(t, expectCode, rec.Code)

    if expectCode == http.StatusOK {
        htmlResponse := rec.Body.String()
        assert.Contains(t, htmlResponse, "listing-row-")
        assert.NotContains(t, htmlResponse, "{\"featured\":")
    }

    if expectCode == http.StatusOK && id == "123" {
        testutil.AssertFeaturedStatus(t, db, id, true)
    }
}
```

**Step 2**: Refactor the test loop block to use the helper, stripping out the `if`s:

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        formData := url.Values{}
        formData.Set("featured", tt.featured)
        urlPath := "/admin/listings/" + tt.id + "/featured"
        if tt.id == "" {
            urlPath = "/admin/listings/featured"
        }
        c, rec := setupAdminTestContext(http.MethodPost, urlPath, strings.NewReader(formData.Encode()))
        setupAdminAuth(t, c)
        if tt.id != "" {
            c.SetParamNames("id")
            c.SetParamValues(tt.id)
        }

        app, h, cleanup := setupAdminTest(t)
        defer cleanup()
        tt.setupData(t, app.DB)

        _ = h.HandleToggleFeatured(c)

        assertFeaturedResponse(t, rec, tt.expectCode, tt.id, app.DB)
    })
}
```

## What NOT to change
- Do not modify the `tests` slice definition
- Do not export the new helper

## Verification

```bash
go test ./internal/module/admin/...
go run cmd/verify/main.go critique 2>&1 | grep "TestAdminHandler_HandleToggleFeatured"
```

The function must no longer appear, or score ≤ 10.

## Commit
```
test(admin): extract assertions to reduce HandleToggleFeatured test complexity
```
