---
description: Phase 2: Autonomous Execution Loop (TDD)
---
# Phase 2: Autonomous Execution Loop (TDD)

**Context Detection**: If a Flash Planning prompt preceded this, skip to Actions.

### Pre-conditions
> **Quota Tripwire**: Is the active model Gemini 3.1 Pro/Opus? [Yes/No]. If Yes, does the prompt contain the exact word `OVERRIDE`? [Yes/No]. If No, HALT immediately and output: *"This is an expensive execution loop. Switch to Gemini 3 Flash or reply OVERRIDE to continue."* Do not proceed until overridden or switched.
- Read `.agents/skills/go-tdd/SKILL.md`.
- Read `.agents/workflows/coding-standards.md` → Testing section.

### Actions
Execute the RED-GREEN-REFACTOR cycle per `.agents/skills/go-tdd/SKILL.md`.

### Verification
- Run `go run ./cmd/verify precommit` before commit to ensure all gates pass.
- **UI Modifications**: If you modified any HTML, CSS, or UI templates, you MUST explicitly run `go run ./cmd/verify browser` locally to catch visual and interaction regressions, as `precommit` intentionally skips them.
