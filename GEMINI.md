## AGENT WORKFLOWS
For detailed rules and development process, run or refer to these workflows:
- `/feature-implementation`: Building out new features with strict TDD.
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
- Transition phases (RED, GREEN, REFACTOR) using `./scripts/agent-exec.sh set-phase <phase>`.
- Verify and pass gates using `./scripts/agent-exec.sh verify <gate_id>`.
- **ANTI-CHEAT**: NEVER manually edit `.agents/state.json`. The file is protected by a cryptographic signature. If you must bypass validation gates for UI/HTML layout changes (where a Go `red-test` is impossible), use `./scripts/agent-exec.sh gate red-test PASS` directly. NOTE: Even if `red-test` is bypassed, you are strictly required to pass the `lint` gate and manually verify the UI changes using the `browser_subagent` before passing the `browser-verification` gate.
- **UPDATE PROGRESS**: Before verifying the `implementation` gate, you MUST create `.tester/tasks/pending_update.json` containing `{"category": "...", "description": "...", "steps": ["..."]}`. The harness will automatically append this to `progress.json` upon Green phase transition.
- **FINALIZE**: After completing a task and all gates are passed, instruct the agent to **ALWAYS** use the `mcp_mcp-memory-service_memory_store` tool to save a summary of the completed feature/bugfix and architectural decisions.