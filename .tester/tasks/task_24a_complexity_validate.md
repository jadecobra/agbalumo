# Task 24a: Reduce Complexity — `(*Listing).Validate` (score: 15 → target ≤ 10)

## File
`internal/domain/listing.go`

## Context

The `Validate()` method scores **15** because it contains two sequential `range` loops
with an `if` inside each, plus an outer `if l.Type == Job` block wrapping a third loop.
Even though the rule slices already exist, the inlined loop bodies add structural nodes.

The existing sub-helpers are already correct:
- `validateOrigin()` — extracted ✅
- `validateTypeRequirements()` — extracted ✅
- `validateContact()` — extracted ✅
- `validateTypeSpecific()` — extracted ✅

**What remains in `Validate()` that drives the score:**

```go
// 2. Simple Field & Length Rules
for _, rule := range validationRules {        // +2 nodes (for + if)
    if rule.condition(l) {
        return errors.New(rule.err)
    }
}
for _, rule := range lengthRules {            // +2 nodes (for + if)
    if rule.field(l) > rule.limit {
        return errors.New(rule.err)
    }
}

// 3. Job Specific Required Fields
if l.Type == Job {                            // +1 node
    for _, f := range jobFields {             // +2 nodes (for + if)
        if f.field(l) == "" {
            return errors.New(f.err)
        }
    }
}
```

## Required Change

Extract a single private helper that iterates all three rule slices:

```go
// applyRules runs the validationRules, lengthRules, and (for Job) jobFields in sequence.
func (l *Listing) applyRules() error {
    for _, rule := range validationRules {
        if rule.condition(l) {
            return errors.New(rule.err)
        }
    }
    for _, rule := range lengthRules {
        if rule.field(l) > rule.limit {
            return errors.New(rule.err)
        }
    }
    if l.Type == Job {
        for _, f := range jobFields {
            if f.field(l) == "" {
                return errors.New(f.err)
            }
        }
    }
    return nil
}
```

Then `Validate()` becomes:

```go
func (l *Listing) Validate() error {
    if err := l.validateOrigin(); err != nil {
        return err
    }
    if err := l.validateTypeRequirements(); err != nil {
        return err
    }
    if err := l.validateContact(); err != nil {
        return err
    }
    if err := l.applyRules(); err != nil {
        return err
    }
    return l.validateTypeSpecific()
}
```

## What NOT to change
- Do not alter any existing sub-helpers (`validateOrigin`, `validateContact`, etc.)
- Do not modify any rule slice definitions
- Do not touch any test files

## Verification

```bash
go test ./internal/domain/...
go run cmd/verify/main.go critique 2>&1 | grep "domain"
```

The `(*Listing).Validate` line must no longer appear in the cognitive complexity output,
OR its score must be ≤ 10.

## Commit

```
refactor(domain): extract applyRules helper to reduce Validate complexity
```
