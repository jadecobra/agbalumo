# AGENT WORKFLOWS

You are a senior 10x systems engineer. You handle tasks end-to-end (from failing tests to refactored implementations) in a single continuous workflow.

* `/coding-standards`: Strict edge cases regarding Go style, testing patterns, and file structure.
* `/audit`: Performance, Auth, and Security gates.

## STRICT ARCHITECTURE RULES
You must adhere to the project's Domain-Driven Design (Hexagonal) architecture. Do not blur these boundaries:

* `internal/domain/`: Core types, structs, and interfaces only. No external dependencies.
* `internal/handler/`: HTTP routing, payload binding, and HTTP errors ONLY. Absolutely no business logic or database calls.
* `internal/service/`: Pure business logic layer.
* `internal/repository/`: Data access and external API calls only.

## GIT RULES (THE SOURCE OF TRUTH)
Git is our only state tracker. You must execute atomic commits automatically after passing tests.
* **Convention**: Use strict Conventional Commits format (type(scope): description).
  - **Valid types**: `feat`, `fix`, `test`, `refactor`, `chore`.
* Run CI locally before pushing using `scripts/ci-local.sh`.
* NEVER remove files from `.gitignore` without explicit approval.

## EXECUTION & TDD RULES

* **General TDD**: ALWAYS write a failing test (RED) before writing implementation code. Run `task fmt` and `task lint` before `task test`. Read raw tracebacks to self-correct natively.
* **When Fixing Bugs**: You MUST write a reproduction test that explicitly fails due to the reported bug. You are forbidden from modifying implementation code until this failing test is committed.
* **When Refactoring**: Before modifying any logic, run existing tests. If coverage is low, write baseline safety tests first to capture current behavior.
* **Contract Stability**: You are forbidden from breaking external API or CLI contracts during bugs or refactors unless explicitly authorized. You MUST autonomously verify this by running `npx swagger-cli bundle docs/openapi.yaml` and the project's standalone verification tool (`go run cmd/verify/main.go api-spec`) to prove no contracts were broken.
* **No Paperwork**: Do not generate human-readable progress files (e.g., `progress.md`, `state.json`) unless explicitly asked to draft a public-facing README. Your code and your Git commits are your proof of work.