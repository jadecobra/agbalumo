---
name: "Flash Review"
description: "Post-Flash execution verification checklist — deterministic action-by-action review."
triggers:
  - "review flash output"
  - "check implementation"
  - "verify flash changes"
  - "flash implemented"
mutating: false
---

# Flash Review Skill

## Trigger
After a Flash model implements changes from a planning prompt, use this skill to verify correctness.

## Procedure

### Phase 1 — Action Verification (Deterministic)
For each action item in the original prompt:
1. Use `grep_search` or `view_file` to confirm the change was applied at the specified location
2. Check: Does the change match the prompt's specification exactly?
3. Record: ✅ Applied / ❌ Missing / ⚠️ Partial

### Phase 2 — Gate Sweep (Deterministic)
Run `go run ./cmd/verify sweep`. Record pass/fail for each gate.

### Phase 3 — TDD Verification
If the prompt included a TDD section:
1. Run the specified test commands
2. Verify all tests pass
3. Check test names match the prompt's specification

### Phase 4 — CI & Push Verification
1. Verify that the agent executed `./scripts/pushw.sh`.
2. Confirm the remote CI run passed (you can manually verify using `gh run watch` if needed).
3. Do NOT mark the review as complete or CLEAN until the remote CI is green.

### Phase 5 — Residual Scan (Reasoning)
1. Search for stale references related to the changes (e.g., old paths, old function names)
2. Check if the prompt's commit message convention was followed
3. Identify any secondary effects the changes might have caused

### Output Format
```
## Flash Review: [Prompt Name]
### Actions: X/Y verified
### Gates: X/Y passed  
### Tests: X/Y passed
### CI Parity: ✅ Passed / ❌ Failed
### Residuals: [count] found
### Verdict: CLEAN / [N] RESIDUALS
```

### Rules
- Do NOT fix residuals during review. Report them for a follow-up prompt.
- If >50% of actions failed verification, flag the entire prompt for re-execution.
- Score delta is optional — only compute if a baseline score was established.
