## AGENT WORKFLOWS
For detailed rules and development process, run or refer to these workflows:
- `/feature-implementation`: Building out new features with strict TDD.
- `/coding-standards`: specific edge cases regarding Go style, testing patterns, and file structure.
- `/audit`: Performance, Auth, and Security gates.
- `/restart-server`: Commands to rebuild CSS and binary.

## Git Rules
- keep commit message short and concise, imperative mood.
- Run CI locally before pushing using `scripts/ci-local.sh`.
- NEVER remove files from `.gitignore` without explicit approval.

## Harness Integration
- **ALWAYS** initialize the 10x Engineer harness when starting a new task (feature, bugfix, or refactor) by running `./scripts/agent-exec.sh init <feature_name> <workflow_type>`.
- Transition phases (RED, GREEN, REFACTOR) using `./scripts/agent-exec.sh set-phase <phase>`.
- Verify and pass gates using `./scripts/agent-exec.sh verify <gate_id>`.
- **UPDATE PROGRESS**: Before verifying the `implementation` gate, you MUST create `.tester/tasks/pending_update.json` containing `{"category": "...", "description": "...", "steps": ["..."]}`. The harness will automatically append this to `progress.json` upon Green phase transition.
- **FINALIZE**: After completing a task and all gates are passed, instruct the agent to **ALWAYS** use the `mcp_mcp-memory-service_memory_store` tool to save a summary of the completed feature/bugfix and architectural decisions.