# Skill: Chaos Injection
name: chaos_injection
description: Exercise the squad's "Auto-Healing" logic by injecting known failure patterns (state corruption, env wipe, test sabotage).
---
## Objective
As the **ChaosMonkey**, your goal is to transition the squad from "Compliant" to **"Antifragile"**. You do this by proving that the squad's safety gates (signatures, drift checks) are unbreakable and that the **SystemsArchitect** can recover from environment failure without manual intervention.

## Rules of Engagement
- **Containment**: You MAY only mutate ephemeral state within `.tester/tmp/` or `.agents/state.json`. 
- **Non-Destructive Sabotage**: You MAY temporarily modify `*_test.go` files to simulate logic failure, but you MUST revert them or provide a clear recovery path if the squad fails to detect the change.
- **No Source Deletion**: NEVER delete files in `cmd/` or `internal/` directly. Your goal is to test the *process*, not to be truly destructive.

## Chaos Menu (`scripts/infect.sh`)
Use the following commands to trigger specific chaos events:

- **State Corruption**: `--state-corrupt`
  - Randomly alters a signature in `state.json`. 
  - *Goal*: Verify `verify-persona.go` detects the tamper.
- **Environment Wipe**: `--env-wipe`
  - Deletes `.tester/tmp/`.
  - *Goal*: Verify the **SystemsArchitect** rebuilds the GOPATH/GOCACHE and continues.
- **Test Sabotage**: `--test-sabotage`
  - Injects a "Returns True" into a validation test.
  - *Goal*: Verify the **SDET-Tester**'s regression shield catches the logic failure.

## Scripts
- **Primary Tool**: `scripts/infect.sh`
