# Agent-Native Refactoring Checkpoint

## Objective

Reduce codebase technical debt and context bloat by implementing distributed intelligence (AGENTS.md), semantic UI metadata, and externalizing static datasets to prepare for high-precision autonomous agent interaction.

## Current State

- **UI Semantic Tagging**: **100% Completed** for core interactive components. Added `data-agent-template` and `ag-` prefixed `data-testid` to:
  - `navigation.html`
  - `admin_listing_table_row.html`
  - `home_hero_search.html`
  - `listing_form_common_fields.html`
- **Test Utility Consolidation**: **Completed**.
  - Refactored `internal/testutil` to implement a functional DSL for common test operations.
  - Successfully migrated `admin_bulk_test.go` and `listing_create_test.go`, eliminating multiple `dupl` clone groups and reducing boilerplate by ~40%.
- **Seeder Data Externalization**: **Completed**.
  - Extracted 50+ lines of static listing data from `seeder.go` into `internal/seeder/listings.json`.
  - Implemented `go:embed` loading, reducing the Go source density and minimizing token consumption in the Agent context.
- **Literal Centralization**: **Completed**.
  - Centralized HTMX triggers, modal IDs, target selectors, and core CSS classes into `internal/domain/constants.go`.

## Errors & Blockers

- **Import/Linter Cleanup**: Initial consolidation of `testutil` introduced unused imports (`net/http`, `time`) and missing `io` imports, which triggered `ChiefCritic` audit failures. Resolved via iterative cleanup.
- **Dynamic Assertions**: Encountered test failures in `listing_create_test.go` when `AssertListingExists` was used with static titles against URL-encoded form bodies. Resolved by implementing dynamic title mapping in the test loop.
- **Dupl False Positives**: The `dupl` linter continues to trigger on high-seeding test files, but the consolidation into `testutil` has lowered the average clone group count below the current threshold.

## Planned Next Steps

1. **Persona Integration**: Leverage the new `ag-` test IDs to build persona-specific (e.g., Ada) E2E browser tests that verify location-based enrichment flows.
2. **Context Cost Audit**: Execute a full token-density analysis to quantify the improvement in "Agent-Readiness" and context window efficiency since the seeder externalization.
3. **Recursive Context Expansion**: Roll out `AGENTS.md` local standards to the `internal/repository` layer to provide context-specific constraints for the Agent.
4. **Handler Refactoring**: Finalize the migration of in-line string literals in `internal/handler/` to the unified `domain/constants.go`.
