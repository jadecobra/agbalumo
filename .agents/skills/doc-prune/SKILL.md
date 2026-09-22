---
name: "Doc Prune"
description: "Audit repository documentation against the 5-tier enforcement model to delete unenforced prose and eliminate documentation drift."
triggers:
  - "/doc-prune"
  - "doc prune"
  - "prune documentation"
  - "prune docs"
mutating: true
---

# Doc Prune Skill

Evaluate documentation files against the Documentation Tier Test.
Delete or fold documents lacking an enforcement mechanism.

## The Documentation Tier Test
Every documentation file must satisfy at least one tier.

1. **Tier 1 (Code)**: Compiler enforcement (Go types, struct definitions, interfaces).
2. **Tier 2 (Verify Tool)**: Deterministic gate enforcement (`verify design`, `verify deprecated`).
3. **Tier 3 (AGENTS.md)**: Session context surfaces the document (`internal/module/AGENTS.md`).
4. **Tier 4 (Strict Lesson)**: Triggered matching in preflight via bracketed tags.
5. **Tier 5 (ADR)**: Architectural decision records in `docs/adr/`.
6. **Tier ∅ (Unenforced Prose)**: Documents lacking any enforcement. Delete or fold into Tier 3.

## Procedure
1. Inventory all markdown files in `docs/` excluding `docs/adr/` and `docs/openapi/`.
2. Classify each document into its appropriate tier.
3. For Tier ∅ files, fold unique technical patterns into the nearest `AGENTS.md` file, then delete the source document.
4. Perform single-writer deduplication across documentation covering similar topics.
5. Run `go run ./cmd/verify doc-drift` to confirm no broken path references exist.
6. Commit changes with `chore(docs): prune Tier ∅ documentation per doc-prune audit`.
