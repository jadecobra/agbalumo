---
description: Audit and complete a skill's infrastructure using the 7-item checklist.
---
# /skillify <skill-name>
Audit a skill against the 7-item completeness checklist and fill gaps.
## Phase 1: Audit
Run both verify commands and note failures:
```bash
go run ./cmd/verify skill-conformance
go run ./cmd/verify check-resolvable
```
Check manually:

SKILL.md — has name, description, triggers[], mutating?
Code — has verify subcommand or procedural steps?
Unit tests — deterministic logic tested?
Resolver entry — listed in .agents/skills/RESOLVER.md?
Conformance — skill-conformance passes?
Resolvability — check-resolvable passes?
Manifest — listed in .agents/verify-manifest.yaml?
Report: "Skill : N/7 complete. Missing: [list]"

Phase 2: Fill Gaps (top-down)
Work items 1→7 in order. Each earlier item constrains later items.

Phase 3: Verify
```bash
go run ./cmd/verify skill-conformance
go run ./cmd/verify check-resolvable
go test ./internal/maintenance/ -run TestSkill
```
Commit: chore(agents): skillify <name> to N/7
