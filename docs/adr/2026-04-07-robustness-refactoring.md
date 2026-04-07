# ADR 002: Robustness Audit Codebase Simplification
**Date**: 2026-04-07 **Status**: Accepted

## Context
The integration of `cmd/verify/main.go critique` identified multiple systemic occurrences across the codebase failing project robustness gates, specifically functions chronically exceeding our cognitive complexity threshold (limit of 10) and sub-optimal memory alignments in structs leading to wasteful padding bytes. Test suite files organically accumulated repeated scaffolding and nested loops, obscuring domain assertions.

## Decision
We systematically remediated test-suite and operational logic cognitive complexity by aggressively extracting setup and looping logic into clearly bounded, single-purpose helper functions (e.g., `generateSingleStressListing`, `verifyValidConfigRun`). Furthermore, we reordered over 28 critical domain and response structs using the strict `fieldalignment` standard to minimize zero-padding bytes in memory. Repetitive clone groups (like UI test renderer initializations and basic auth setup) across tests have also been coalesced into `internal/testutil/`.

## Consequences
The primary codebase and its corresponding test suites now exhibit an optimal memory layout and flattened cognitive graphs, passing `gocognit`, `fieldalignment`, and `dupl` verifications uniformly. Engineering maintainability is significantly eased. However, creating new tests requires developers to be mindful of centralized standard helpers located in `internal/testutil/` to prevent regressions back into copy-paste paradigms, and future structs must strictly honor descending size memory layouts to pass the continuous build step.
