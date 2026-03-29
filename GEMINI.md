## AGENT WORKFLOWS
For detailed rules and development process, run or refer to these workflows:
- `/build-feature`: Building out new features with strict TDD.
- `/harness`: Comprehensive command reference for the execution harness.
- `/coding-standards`: specific edge cases regarding Go style, testing patterns, and file structure.
- `/audit`: Performance, Auth, and Security gates.
- `/restart-server`: Commands to rebuild CSS and binary.

## Git Rules
- **ALWAYS execute `git commit` automatically (using `run_command`) after all gates pass, without waiting for the user to explicitly tell you to do so.**
- keep commit message short and concise, imperative mood.
- Run CI locally before pushing using `scripts/ci-local.sh`.
- NEVER remove files from `.gitignore` without explicit approval.

## Harness Integration
- **ALWAYS** initialize the 10x Engineer harness when starting a new task (feature, bugfix, or refactor) by running `./scripts/agent-exec.sh init <feature_name> <workflow_type>`.
- **MANDATORY REFERENCE**: Refer to the `/harness` workflow for all available commands, gate IDs, and phase transition logic.
- **CONTEXT PROTOCOL**: Before starting any task in the `GREEN` or `REFACTOR` phases, you **MUST** locate and read the active `implementation_plan.md` (linked in `progress.json`) to verify the technical details (thresholds, contracts, and patterns) of the feature.
- **FEEDBACK LOOP PROTOCOL**: Catch errors immediately. **ALWAYS** run `task fmt` and `task lint` *before* attempting `task test`. If automated formatting or linting fails, fix syntax errors *immediately* before proceeding to tests.
- Transition phases (RED, GREEN, REFACTOR) using `./scripts/agent-exec.sh set-phase <phase>`.
- Verify and pass gates using `./scripts/agent-exec.sh verify <gate_id>`.
- **CLI DRIFT**: If you are unsure of subcommand arguments (e.g., `init` vs `start`), ALWAYS run `./scripts/agent-exec.sh --help` to verify usage.
- **ANTI-CHEAT**: NEVER manually edit `.agents/state.json`. The file is protected by a cryptographic signature. Note: manual bypassing of `red-test` via `gate red-test PASS` is explicitly blocked. If you must bypass validation gates for UI/HTML layout changes, use `./scripts/agent-exec.sh verify red-test ui-bypass` directly. **WARNING**: The **ChaosMonkey** persona actively tests these gates via fault injection; any successful bypass without a rejection is a system failure. 
- **UPDATE PROGRESS**: Before verifying the `implementation` gate, you MUST create `.tester/tasks/pending_update.json` containing `{"category": "...", "description": "...", "steps": ["..."]}`. The harness will automatically append this to `progress.json` upon Green phase transition.
- **FINALIZE**: After completing a task and all gates are passed, instruct the agent to **ALWAYS** use the `mcp_mcp-memory-service_memory_store` tool to save a **"Squad-Decision-Summary"** including architectural decisions made by the **SystemsArchitect** and product strategy from the **ProductOwner**.