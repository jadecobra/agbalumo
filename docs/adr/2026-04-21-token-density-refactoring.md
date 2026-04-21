# ADR [005]: Token Density Refactoring
**Date**: 2026-04-21 **Status**: Proposed

## Context
Monolithic utility files (e.g. `cmd/verify/main.go`) exceeding 500 lines/tokens degrade AI Agent performance by inflating context windows with irrelevant logic during pinpoint edits. The existing `cmd/verify/main.go` has grown to over 600 lines, containing a mix of CI logic, contract validation, background jobs, and generic utilities.

## Decision
We will split `cmd/verify/main.go` into smaller, cohesive domain-driven files:
- `cmd/verify/main.go`: Boilerplate and core bootstrapping.
- `cmd/verify/ci.go`: CI orchestration and validation flows (test, precommit, trivy).
- `cmd/verify/drift.go`: Contract and template drift detection.
- `cmd/verify/jobs.go`: Background worker triggers (backfill, enrich).
- `cmd/verify/misc.go`: Audit, coverage, and performance utilities.

## Consequences
- **Easier Agentic manipulation**: Agents will load only relevant domain files, reducing token costs and hallucination risk.
- **Improved Maintainability**: Logical grouping makes finding specific command implementations faster for humans.
- **Project Structure**: Slightly higher file count in `cmd/verify/`, but fundamentally cleaner boundaries.
- **Initialization Fix**: We will consolidate the redundant `init()` calls that currently cause duplicated command registration.
