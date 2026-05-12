---
description: Lightweight refactoring workflow with proactive improvement scan.
---

# /refactor <target>

Lightweight alternative to `/build-feature` for pure refactoring tasks (renaming, structural moves, dead code removal, pattern migration). Skips product interrogation.

## Phase 1: Baseline (Deterministic)
> **Quota Tripwire**: Is the active model Gemini 3.1 Pro/Opus? [Yes/No]. If Yes, does the prompt contain the exact word `OVERRIDE`? [Yes/No]. If No, HALT immediately and output: *"Refactoring is an expensive execution loop. Switch to Gemini 3 Flash or reply OVERRIDE to continue."* Do not proceed until overridden or switched.
1. Run `go run ./cmd/verify preflight`
2. Run `go test ./...` — record pass count and coverage as baseline
3. Run `go run ./cmd/verify deprecated` — record violation count as baseline
4. Run `go run ./cmd/verify critique` — record any existing issues

## Phase 2: Execute Refactor (TDD)
**Pre-condition**: Read `.agents/skills/go-tdd/SKILL.md`

1. If the refactor changes behavior, write a test capturing current behavior FIRST
2. Execute the refactor
3. Run `go test ./...` — all tests must still pass
4. Run `go run ./cmd/verify deprecated` — violation count must not increase
5. Commit: `refactor(<scope>): <description>`

## Phase 3: Proactive Improvement Scan — "How Can I Make This Better?"

After the refactor passes all gates, pause and ask:

> **"The refactor is complete. Before closing, let me scan for non-obvious improvements in the affected area."**

Run these tools on the **modified packages only**:
```bash
go run ./cmd/verify context-cost          # Are any files too large now?
go run ./cmd/verify deprecated            # Did the refactor expose new migration opportunities?
go run ./cmd/verify agents-coverage       # Did we create new packages without AGENTS.md?
go run ./cmd/verify critique --baseline HEAD~1  # Did complexity improve or regress?
```

Then answer these questions:
1. **Structural**: Could any function/file be split for better single-responsibility?
2. **Naming**: Are there any names that no longer reflect their purpose after the refactor?
3. **Dead code**: Did the refactor orphan any helpers, types, or constants?
4. **Test gaps**: Are there untested edge cases in the refactored code?
5. **Patterns**: Does the refactored code follow the same patterns as peer code in the package?

Present ≤3 **tool-grounded** improvement suggestions to the user. Each suggestion must reference a specific file:line and a specific verify tool output.

**Rule**: Suggestions are proposals, not actions. WAIT for user approval before implementing any improvement. The user may say "ship it as-is" — that is a valid response.

## Phase 4: Finalize
1. Run `go run ./cmd/verify precommit`
2. Compare test count and coverage against Phase 1 baseline — must not decrease
3. If improvements were approved and implemented, amend the commit or add a follow-up commit
4. **Remote Verification**: Execute `./scripts/pushw.sh` or `gh run watch` immediately after pushing to ensure the remote CI remains green. Work is not complete until remote verification passes.

## Completion
Output a single-line summary: `refactor(<scope>): <what changed> | tests: <pass>/<total> | deprecated: <before>→<after>`
