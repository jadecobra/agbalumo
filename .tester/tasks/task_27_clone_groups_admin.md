# Task 27: Clone Group Reduction — Batch B (admin module production code)

## Context

The `internal/module/admin/admin.go` file contains the highest-impact production
code clone groups. This task targets them specifically, scoped to the admin module's
non-test files.

---

## Target Clone Groups

### 1. `admin.go:136–170` — 6-clone `errgroup.Go` wrapper pattern

**Pattern**: Every `g.Go(func() error { var err error; X, err = h.App.DB.Y(ctx); return err })` 
block in `HandleDashboard` is identical in structure — only the variable assigned
and the DB method called differ.

**Fix**: There is no clean generic abstraction in pure Go without reflection. Instead,
consolidate into a **single sequential load function** `loadDashboardData(...)` that
returns all values and one error. Accept the trade-off: lose micro-parallelism in the
dashboard, gain clarity and compliance.

> Alternatively: define typed adapter helpers like:
> ```go
> func fetchClaims(ctx context.Context, db ...) func() error {
>     return func() error { var err error; result, err = ...; return err }
> }
> ```
> and pass `result` by pointer. Either approach is acceptable.

### 2. `admin.go:69–71` and `admin.go:106–108` — 3-clone redirect-if-no-user blocks

**Pattern**:
```go
u, ok := user.GetUser(c)
if !ok || u == nil {
    return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
}
```
This appears in `AdminMiddleware`, `HandleLoginAction`, and `listing_profile.go`.

**Fix**: Extract `requireUser(c echo.Context) (*domain.User, error)` helper in the
`user` or `admin` package. Callers check the error and return it directly.

### 3. `admin.go:91` and `admin.go:99` — 2-clone `admin_login.html` render calls

**Pattern**: `c.Render(http.StatusOK, "admin_login.html", map[string]interface{}{...})`
appears twice in the same function with nearly identical payloads.

**Fix**: Extract `renderLoginView(c echo.Context, errMsg string) error` helper.
Empty `errMsg` renders without the Error key; non-empty includes it.

### 4. `admin_listings.go:70` and `admin_listings.go:99` — 2-clone `admin_listing_table_row` template calls

**Pattern**: `c.Render(http.StatusOK, "admin_listing_table_row", ...)` repeated
with similar data shapes.

**Fix**: Extract `renderListingRow(c echo.Context, listing domain.Listing) error`.

---

## Verification

- [ ] `go test ./internal/module/admin/...`
- [ ] `go run cmd/verify/main.go critique` — confirm admin clone count is reduced
- [ ] Commit: `refactor(admin): extract shared handler helpers to reduce clone groups`
