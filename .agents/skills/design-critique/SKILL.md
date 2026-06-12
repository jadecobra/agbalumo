---
name: "Design Critique"
description: "3-phase design critique — deterministic audit, Flash grading, optional taste review."
triggers: ["/design-critique", "critique design", "review ui", "harsh review"]
mutating: false
---
# /design-critique (v3)

**Phase 1 — Deterministic Audit** (zero model cost):
1. Run `go run ./cmd/verify design-evidence --json > /tmp/design-evidence.json`
2. Run `go run ./cmd/verify browser`
3. If violations exist, output them. These are facts, not opinions.

**Phase 2 — Flash Grading** (cheap model):
1. Read the JSON report from Phase 1 (`/tmp/design-evidence.json`).
2. Score each of the 7 dimensions (0-10) using the rubric in `./rubric.md`.
3. For each violation, generate a grep-anchored fix prompt following `.agents/skills/flash-plan/SKILL.md`.
4. Output: Scored critique artifact + fix prompts.

**Phase 2.5 — Rubric Validation Gate** (mandatory before emitting fix prompts):
1. For each proposed fix, cite the EXACT rubric line and point deduction it resolves.
2. If a fix cannot be mapped to a specific `-N` deduction line, label it "Best Practice (non-scoring)" and omit it from the fix prompt list.
3. Use only verified tool outputs as evidence. Do not infer template behavior from static source — conditional rendering (`{{ if }}/{{ else }}`) cannot be verified without execution.

**Phase 3 — Taste Review** (expensive model — ONLY if Phase 2 score < 7.0 OR user requests):
1. View 2-3 screenshots (browser_subagent, minimal interaction).
2. Apply the Subtract Mandate — identify ONE element to delete.
3. Override or adjust Phase 2 scores with visual judgment.
4. Output: Final scores + subtract targets.

---

**Pre-Edit Hunk Plan** (mandatory before any file mutation):
- List ALL hunks to be changed across ALL files.
- For each file with >1 hunk, confirm they will be batched in a single `multi_replace_file_content` call.
- Only after the plan is complete, execute edits.
