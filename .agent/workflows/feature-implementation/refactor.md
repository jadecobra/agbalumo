### 2b. Refactor & Regression
- *Run full regression*:
  ```bash
  // turbo
  go test -race ./...
  ```
- **Gate**: Zero regressions.

### 2c. Pre-Commit Quality Gate
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
