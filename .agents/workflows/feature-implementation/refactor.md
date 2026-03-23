### 2b. Evaluation & Optimization (Refactor Phase)
**Goal**: The "refactor" phase is the opportunity to make the codebase better and find abstractions, instead of waiting to be prompted. Take this time to explore the solution space holding the codebase securely with tests.

- **Instructions**:
  1. *Expand Search*: Look beyond the lines you just added. Evaluate the broader package/component.
  2. *Code Smells*: Actively eliminate duplication, remove hardcoded values, and reduce cyclic complexity (e.g. nested `if-else`, prefer guard clauses).
  3. *Performance & Idioms*: Check for N+1 DB queries and memory allocation overhead in loops. Follow `@[.agent/workflows/coding-standards.md]` for idiomatic Go guidelines.
  4. *Modularity*: Determine if new logic can be refactored into "lego bricks" (helper functions or package-level utilities) that simplify the calling code.

- **Refactoring Constraints & Autonomy**:
  - The agent has full autonomy to restructure the implementation as long as tests pass.
  - **CRITICAL**: The agent MUST NEVER modify existing tests in this stage. It can only *add* tests to cover new abstractions. If an existing test breaks, the refactor is invalid and must be reverted/fixed quickly.
  - **Commit Frequency**: Commit early and frequently (`git commit`) after small successful refactoring loops. Rollbacks should be cheap and feedback loops immediate.

### 2c. Refactor Regression Check
- *Run full regression*:
  ```bash
  // turbo
  go test -race ./...
  ```
- **Gate**: Zero regressions.

### 2d. Pre-Commit Quality Gate
- *Run*:
  ```bash
  ./scripts/pre-commit.sh
  ```
- **Gate: `lint`**
  - **PASS**: Output contains `✅ GolangCI-Lint passed` and exit code is 0.
  - **FAIL**: Lint errors found or config invalid.

- **Gate: `coverage`**
  - **PASS**: Output contains `Coverage: [X]%` where `X >= threshold` from `@[.agent/coverage-threshold]`.
  - **FAIL**: Coverage is below the required threshold.
