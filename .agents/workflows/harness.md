---
description: Comprehensive command reference and protocol guide for the Agbalumo execution harness.
---
Use this workflow to understand how to interact with the project's strictly governed development harness (`scripts/agent-exec.sh`).

### Core Commands

Every command should be run via the wrapper script to ensure correct environment isolation and default formatting:

1. **Initialize a Task**:
   ```bash
   ./scripts/agent-exec.sh init <feature_name> [workflow_type]
   ```
   *   `workflow_type`: `feature` (default), `bugfix`, or `refactor`.
   *   *Note: This creates the `.agents/state.json` file which tracks your progress.*

2. **Transition Phases**:
   ```bash
   ./scripts/agent-exec.sh set-phase <PHASE>
   ```
   *   `PHASE`: `RED` (Writing tests), `GREEN` (Implementation), `REFACTOR` (Optimization), `IDLE`.

3. **Verify Gates**:
   ```bash
   ./scripts/agent-exec.sh verify <gate_id> [pattern]
   ```
   *   `gate_id` options:
       *   `red-test`: Run failing tests. (Use `ui-bypass` as pattern for pure UI changes).
       *   `api-spec`: Check for API/Contract drift.
       *   `implementation`: Run implementation tests (Logic).
       *   `lint`: Run `golangci-lint`.
       *   `coverage`: Verify test coverage thresholds.
       *   `template-drift`: Ensure HTML templates match domain models.
       *   `security-static`: Run AST-based security scanners.
       *   `browser-verification`: Manually mark UI verification as passed.

4. **Check Status**:
   ```bash
   ./scripts/agent-exec.sh status
   ```

### 10x Engineer Protocol Rules

- **Strict Sequence**: You cannot verify `implementation` until `red-test` and `api-spec` are `PASS`.
- **Auto-Transition**: Successfully passing all gates in a phase will automatically transition the state to the next phase.
- **Progress Tracking**: Before verifying `implementation`, you MUST create `.tester/tasks/pending_update.md` to update the project's global `progress.md`.
- **Unsure?**: If you are uncertain about subcommands or arguments, ALWAYS run `./scripts/agent-exec.sh --help`.
