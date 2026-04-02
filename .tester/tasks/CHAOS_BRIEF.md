---
target_persona: "@ChaosMonkey"
phase: CHAOS
feature: upgrade-ci-node24
entry_condition: "SecurityEngineer has passed audit_security gate"
---

# Chaos Brief: upgrade-ci-node24

> [!IMPORTANT]
> This file is authored by ProductOwner during Phase 1 **before** autonomous execution begins.
> ChaosMonkey reads this file on `/resume` instead of a generic HANDOFF.

## Authorized Sabotage Targets

Replace these placeholders with feature-specific targets before handing off to ChaosMonkey:

- [ ] **Runner Runtime Incompatibility**: Override the runner's default behavior to use Node 14 (pre-deprecation) — verify the workflow fails with a clear message about Node 24 requirement if `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24` is true.
- [ ] **Action SHA Mismatch**: Intentionally change `actions/checkout` or `actions/cache` SHA to an invalid string in `ci.yml` — verify GitHub Actions reports "Loading action" failure.
- [ ] **Environment Flag Interference**: Set `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: false` in `ci.yml` — verify `@ChiefCritic` or `@SecurityEngineer` detects the regression (the warning returns).
- [ ] **Node Version Regression**: Revert `node-version: '24'` to `'20'` in one of the jobs — verify the deprecation warning reappears in the CI annotations.

## Squad Success Condition

The squad **passes** the chaos test if:
- Any intentional version regression or flag mismatch is detected by the CI or the `@ChiefCritic` audit.
- GitHub's own security gates (SHA verification) prevent execution of tampered action SHAs.
- The `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24` flag successfully suppresses warnings for compatible legacy actions.

## Squad Failure Condition

The squad **fails** if ChaosMonkey successfully regresses the Node version without the CI or auditors flagging the deprecation warning.

On failure: ChaosMonkey kicks back to BackendEngineer (via SDET-Tester) to patch the gap.

## ChaosMonkey Entry Instruction

When you open your chat window and run `/resume`:
1. Read this file first (it is your mission brief, not the generic HANDOFF).
2. Confirm entry condition: `./scripts/agent-exec.sh status --text` must show `security-static: PASSED`.
3. Execute sabotage targets in order.
4. Report results against success/failure conditions above.
5. **Never write fixes.** Kick back to SDET-Tester for any failure mode detected.
