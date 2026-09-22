# ADR 016: Migrate Legacy Workflows to Consolidated Modern Skills
**Date**: 2026-09-22 **Status**: Accepted

## 1. Context & User Problem
Legacy workflow files in `.agents/workflows/` lacked first-class slash command registration, multi-file references, and structured frontmatter schemas.
Multiple workflow files contained redundant steps that duplicated existing core skills (`go-tdd`, `ci-parity`, `verify-authoring`).
Fourteen fragmented workflow files increased context token usage and slowed agent initialization.

## 2. Decision
Migrate deprecated workflows into modern skills under `.agents/skills/<name>/SKILL.md`.
Consolidate redundant workflows rather than creating one-to-one legacy clones.
Fold `build-feature` phases 1, 2, and 3 into a single unified `build-feature` skill.
Integrate `hotfix`, `refactor`, and `debug` triage into `go-tdd`.
Integrate the 7-item skill completeness checklist into `verify-authoring`.
Migrate `learn`, `red-team`, `doc-prune`, `deploy-secrets`, and `stress-test` as independent modern skills.
Reclassify `coding-standards.md` from a pseudo-workflow into canonical repository standards at `.agents/coding-standards.md`.
Completely delete the `.agents/workflows/` directory and purge all `.bak` files.
Update maintenance tools (`preflight`, `lessons`, `session_context`, `minify_context`) to resolve `.agents/coding-standards.md` directly with zero legacy fallback code.

## 3. The Complexity Kill-Switch (Rationale)
* **User Value**: The user retains every slash command with faster command discovery, zero duplicate steps, and unified guidance.
* **Performance Budget**: Reduces total context token footprint by pruning 8 redundant workflow files and collapsing 3 multi-phase files.
* **Minimalism Check**: Completely removed the `.agents/workflows/` directory and eliminated backwards-compatibility shim logic.

## 4. Consequences
* **Technical Tradeoffs**: Skill directories require strict YAML frontmatter (`name`, `description`, `triggers`, `mutating`) enforced by `verify skill-conformance`.
* **Observability**: Monitored via `verify check-resolvable`, `verify skill-conformance`, and `verify doc-drift`.
* **SQLite Impact**: Zero database impact.

## 5. Alternatives Considered
* Direct 1-to-1 migration of all 14 files into 14 skills. Rejected because it preserves dead weight and creates single-caller pseudo-skills.
* Maintaining backwards-compatible symlinks and Go fallback paths. Rejected because maintaining legacy compatibility layers preserves tech debt when modern skills are the active standard.
