## AGENT WORKFLOWS
For detailed rules and development process, run or refer to these workflows:
- `/feature-implementation`: Building out new features with strict TDD.
- `/coding-standards`: specific edge cases regarding Go style, testing patterns, and file structure.
- `/audit`: Performance, Auth, and Security gates.
- `/restart-server`: Commands to rebuild CSS and binary.

## Git Rules
- keep commit message short and concise, imperative mood.
- NEVER commit `ARCHITECTURE_CRITIQUE.md` or remove it from `.gitignore`.
- Run CI locally before pushing using `scripts/ci-local.sh`.
- NEVER remove files from `.gitignore` without explicit approval.