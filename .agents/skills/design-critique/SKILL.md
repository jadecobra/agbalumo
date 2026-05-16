---
name: "Design Critique"
description: "3-phase design critique — deterministic audit, Flash grading, optional taste review."
triggers: ["/design-critique", "critique design", "review ui", "harsh review"]
mutating: false
---
# /design-critique (v2)

**Phase 1 — Deterministic Audit** (zero model cost):
1. Run `go run ./cmd/verify visual-audit --json > /tmp/visual-audit.json`
2. Run `go run ./cmd/verify browser`
3. If violations exist, output them. These are facts, not opinions.

**Phase 2 — Flash Grading** (cheap model):
1. Read the JSON report from Phase 1.
2. Score each of the 6 dimensions (0-10) using the rubric in `./rubric.md`.
3. For each violation, generate a grep-anchored fix prompt following `.agents/skills/flash-plan/SKILL.md`.
4. Output: Scored critique artifact + fix prompts.

**Phase 3 — Taste Review** (expensive model — ONLY if Phase 2 score < 7.0 OR user requests):
1. View 2-3 screenshots (browser_subagent, minimal interaction).
2. Apply the Subtract Mandate — identify ONE element to delete.
3. Override or adjust Phase 2 scores with visual judgment.
4. Output: Final scores + subtract targets.
