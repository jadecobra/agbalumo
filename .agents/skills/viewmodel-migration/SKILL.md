---
name: "ViewModel Migration"
description: "Migrate handlers from map[string]interface{} to typed ViewModel structs."
triggers:
  - "migrate handler"
  - "typed viewmodel"
  - "fix deprecated map"
  - "viewmodel migration"
mutating: true
---

# ViewModel Migration Skill

## Trigger
When `verify deprecated` reports `map[string]interface{}` violations in a handler.

## Pre-conditions
- Run `go run ./cmd/verify deprecated` to identify target files
- Read `internal/module/base.go` for `BaseViewData` struct and `PopulateBase` pattern
- Read `internal/module/AGENTS.md` for Quick Reference

## Procedure

### Step 1: Define ViewModel
In the handler's package, create a struct:
```go
type XxxViewModel struct {
    module.BaseViewData
    // Add fields matching template {{ .Field }} references
    SpecificField string
    Items         []domain.Item
}
```

### Step 2: Identify Template Fields
Run `grep -n '{{ \.' ui/templates/<template>.html` to find all field references.
Cross-reference against the new struct. Every `{{ .Field }}` must resolve.

### Step 3: Refactor Handler
```go
// BEFORE (deprecated)
data := map[string]interface{}{
    "Items": items,
}
return h.RenderWithBaseContext(c, "template", data)

// AFTER (typed)
vm := XxxViewModel{
    BaseViewData:  h.PopulateBase(c),
    Items:         items,
}
return h.RenderTyped(c, "template", vm)
```

### Step 4: Verify
```bash
go run ./cmd/verify template-contract   # Fields resolve
go run ./cmd/verify deprecated          # Violation count decreased
go test ./internal/module/<pkg>/ -v     # No regressions
```

### Step 5: Commit
`refactor(<scope>): migrate <handler> to typed ViewModel`

## Anti-patterns
- Do NOT use `interface{}` for any ViewModel field — use concrete types
- Do NOT embed data in the template via `dict` helper if a struct field works
- Do NOT skip `PopulateBase` — it handles Categories, User, CSRF, Env
