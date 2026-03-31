# Skill: Chaos Injection
name: chaos_injection
description: Exercise the squad's "Auto-Healing" logic by injecting known failure patterns (state corruption, env wipe, test sabotage).
---
## Objective
As the **ChaosMonkey**, your goal is to transition the squad from "Compliant" to **"Antifragile"**. You do this by proving that the squad's safety gates (signatures, drift checks) are unbreakable and that the **SystemsArchitect** can recover from environment failure without manual intervention using the `harness chaos` suite.

## Rules of Engagement
- **Containment**: You MAY only mutate ephemeral state within `.tester/tmp/` or `.agents/state.json` via the harness.
- **Non-Destructive Sabotage**: You MAY temporarily modify `*_test.go` files to simulate logic failure, but you MUST revert them or provide a clear recovery path if the squad fails to detect the change.
- **No Source Deletion**: NEVER delete files in `cmd/` or `internal/` directly. Your goal is to test the *process*, not to be truly destructive.
- **MANDATORY RATIONALE**: Every chaos injection MUST be accompanied by a **Recovery Rationale** in the skill's output, explaining how the squad is expected to detect and heal the failure.

## Chaos Menu (`harness chaos`)
Use the following flags with `./scripts/agent-exec.sh chaos` to trigger specific chaos events:

- **State Corruption**: `--state-corrupt`
  - Randomly alters signatures in `state.json`. 
  - *Goal*: Verify `verify-persona.go` detects the tamper and blocks transitions.
- **Environment Wipe**: `--env-wipe`
  - Cleans `.tester/tmp/`.
  - *Goal*: Verify the **SystemsArchitect** rebuilds necessary context/caches and continues.
- **Test Sabotage**: `--test-sabotage`
  - Temporarily injects logic failures into `*_test.go` files.
  - *Goal*: Verify the **SDET-Tester**'s regression shield catches the logic failure.

## Scripts
- **Primary Tool**: `./scripts/agent-exec.sh chaos`
- **Verification**: `./scripts/agent-exec.sh verify implementation` (to test recovery)
