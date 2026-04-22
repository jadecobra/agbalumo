# ADR: Domain Listing Model Split

## Context
`internal/domain/listing.go` has grown into a monolithic file (282 lines) that mixes core data structures (`Listing`), massive hardcoded country mappings (`ValidOrigins`), and heavy validation logic. This bloats the token context for Agents when they only need to modify the data structure, leading to potential inaccuracies and inefficiency.

## Decision
We will split `internal/domain/listing.go` into three focused files within the same `domain` package:
1. `listing.go`: Core struct and constant definitions.
2. `listing_validation.go`: All validation rules and method implementations (`Validate`, `applyRules`, etc.).
3. `origins.go`: The `ValidOrigins` map and origin-specific validation.

## Consequences
- **Reduced Token Context**: Agents working on the `Listing` struct no longer pull in 200+ lines of validation and country lists.
- **Improved Maintainability**: Validation logic is isolated, making it easier to audit and extend without touching the core model.
- **No API Breakage**: Since all files remain in the `domain` package, there are no changes to external call sites.
- **Cleaner Diffing**: Changes to validation rules will no longer clutter the history of the core model file.
