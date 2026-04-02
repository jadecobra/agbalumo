---
target_persona: "@ChaosMonkey"
phase: CHAOS
feature: REPLACE_WITH_FEATURE_NAME
entry_condition: "SecurityEngineer has passed audit_security gate"
---

# Chaos Brief: REPLACE_WITH_FEATURE_NAME

> [!IMPORTANT]
> This file is authored by ProductOwner during Phase 1 **before** autonomous execution begins.
> ChaosMonkey reads this file on `/resume` instead of a generic HANDOFF.

## Authorized Sabotage Targets

Replace these placeholders with feature-specific targets before handing off to ChaosMonkey:

- [ ] **State Corruption**: Corrupt `.agents/state.json` after GREEN phase passes — verify ANTI-CHEAT triggers.
- [ ] **Gate Bypass Attempt**: Attempt `./scripts/agent-exec.sh gate coverage PASS` — verify it is blocked and logged in `bypass_audit.log`.
- [ ] **Test Sabotage**: Introduce a subtle off-by-one or nil pointer in `REPLACE_WITH_TEST_FILE` — verify SDET-Tester detects it within 1 regression run.
- [ ] **HANDOFF Staleness**: Replace `HANDOFF.md` with a stale file from a prior feature — verify Step 1.5 tripwire triggers RESUME HALTED.

## Squad Success Condition

The squad **passes** the chaos test if:
- SDET-Tester detects every sabotage via regression test within 1 attempt.
- The harness ANTI-CHEAT fires on state corruption.
- The bypass audit log contains an entry for the blocked coverage gate attempt.

## Squad Failure Condition

The squad **fails** if ChaosMonkey successfully bypasses any gate or corrupts state without the system detecting it.

On failure: ChaosMonkey kicks back to BackendEngineer (via SDET-Tester) to patch the gap.

## ChaosMonkey Entry Instruction

When you open your chat window and run `/resume`:
1. Read this file first (it is your mission brief, not the generic HANDOFF).
2. Confirm entry condition: `./scripts/agent-exec.sh status --text` must show `security-static: PASSED`.
3. Execute sabotage targets in order.
4. Report results against success/failure conditions above.
5. **Never write fixes.** Kick back to SDET-Tester for any failure mode detected.
