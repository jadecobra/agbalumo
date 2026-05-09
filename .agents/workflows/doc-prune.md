---
description: Documentation tier test — evaluate and prune docs that lack enforcement mechanisms.
---

# /doc-prune

Evaluate all documentation files against the Documentation Tier Test. Delete or fold docs that have no enforcement mechanism.

## The Tier Test

Every piece of agent-facing documentation must satisfy **at least one** tier:

| Tier | Mechanism | Enforcement | Example |
|------|-----------|-------------|---------|
| 1 | Code | Go compiler | `BaseViewData` struct |
| 2 | Verify tool | `precommit` / `ci` gates | `verify design`, `verify deprecated` |
| 3 | AGENTS.md | `session-context` surfaces it | `internal/module/AGENTS.md` |
| 4 | Strict Lesson | `preflight` matches triggers | `[TRIGGER: handler_change]` |
| 5 | ADR | `session-context` matches | `docs/adr/2026-05-05-*.md` |
| ∅ | Standalone prose | **Nothing enforces it** | ← Delete or fold into Tier 3 |

## Procedure

### Step 1: Inventory
List all `.md` files in `docs/` (excluding `adr/` and `openapi/`):
```bash
find docs/ -name "*.md" -not -path "docs/adr/*" -not -path "docs/openapi/*"
```

### Step 2: Classify
For each file, determine its tier:
- Does a `verify` tool enforce it? → Tier 2
- Is it referenced by an AGENTS.md `session-context`? → Tier 3
- Does it contain `[TRIGGER:]` tagged lessons? → Tier 4
- Is it an ADR? → Tier 5
- None of the above? → Tier ∅ (candidate for deletion)

### Step 3: Act on Tier ∅ Files
For each Tier ∅ file:
1. Check if it contains unique, useful content (function signatures, patterns, constraints)
2. If YES: fold the useful content into the nearest AGENTS.md (Tier 3) and delete the original
3. If NO: delete the file outright
4. Run `go run ./cmd/verify doc-drift` to ensure no stale references remain

### Step 4: Single Writer Audit
Check for duplicate semantic coverage:
- Two files covering "coding standards"? → Keep the enforced one, delete the other
- A doc repeating what AGENTS.md already says? → Delete the doc

### Output
```
## Doc Prune Results
| File | Tier | Action | Reason |
|------|------|--------|--------|
```

### Commit
`chore(docs): prune Tier ∅ documentation per doc-prune audit`
