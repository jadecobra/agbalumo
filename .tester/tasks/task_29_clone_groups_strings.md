# Task 29: Clone Group Reduction — Batch D (string constants & repeated literals)

## Context

The `[2/4] Repeated Strings` check does not currently fail the gate, but the
flagged literals are also contributors to clone groups. This batch extracts the
most impactful repeated string literals into package-level constants, closing both
the strings check and overlapping clone groups simultaneously.

This is the **lowest-risk, most mechanical** task in the series. Flash can handle
the full batch in one turn.

---

## Target Repeated Literals

### 1. `internal/repository/sqlite/` — `"listing not found"` (2 locations)

Files: `sqlite_listing_write.go:146`, `sqlite_listing_read.go:169`

**Fix**: Add to `internal/repository/sqlite/queries.go` (or a new `errors.go`):
```go
const errListingNotFound = "listing not found"
```
Use in both files.

### 2. `internal/repository/sqlite/sqlite_claim.go` — `"claim request not found"` (2 locations)

Lines 62 and 94.

**Fix**: 
```go
const errClaimNotFound = "claim request not found"
```

### 3. `internal/module/listing/listing_mutations.go` — `"Title already exists..."` (2 locations)

Lines 42 and 119. Already a long string — extract:
```go
const errTitleExists = "Title already exists. Please choose a different title."
```

### 4. `internal/module/feedback/feedback.go` and `listing_mutations.go` — `"Login required"` (2 locations)

Extract to a shared constants file or `internal/common/errors.go`:
```go
const ErrMsgLoginRequired = "Login required"
```

### 5. `internal/module/admin/admin_claims.go` — `"Claim request not found"` (2 locations)

Lines 16 and 28.

**Fix**: Package-level const in `admin_claims.go`:
```go
const errClaimRequestNotFound = "Claim request not found"
```

### 6. `internal/module/listing/listing_form.go` — `"2006-01-02T15:04"` datetime format (2 locations)

Lines 96 and 103.

**Fix**:
```go
const datetimeLocalFormat = "2006-01-02T15:04"
```

### 7. `internal/module/admin/admin_listings.go` — `"admin_listing_table_row"` (2 locations)

Lines 70 and 99.

**Fix**:
```go
const tmplListingTableRow = "admin_listing_table_row"
```

### 8. `internal/infra/server/server.go` and `internal/module/auth/provider.go` — `"BASE_URL"` (2 locations)

**Fix**: Move to `internal/config/config.go` as:
```go
const EnvBaseURL = "BASE_URL"
```

### 9. `internal/maintenance/perf.go` — `"could not find %s: %w"` (2 locations)

Lines 64 and 72.

**Fix**:
```go
const errNotFoundFmt = "could not find %s: %w"
```

---

## Verification

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go run cmd/verify/main.go critique` — confirm string clone count reduction in affected packages
- [ ] Commit: `refactor(constants): extract repeated string literals into named constants`
