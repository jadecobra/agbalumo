### 2a. Implement Logic (GREEN)

> **Persona: Backend** — Write only the minimal code to pass existing tests. Validate all inputs (trust nothing). Use Echo framework. Maintain minimal external dependencies. Run benchmarks for critical logic (`go test -bench=.`, target < 1000ns/op).

- Write the logic in `handler/`, `service/`, etc.
- Also ensure tests cover `400 Bad Request` schema validations and `401/403` auth checks defined in the spec.
- *Run*:
  ```bash
  // turbo
  go test -v -run TestNewFeatureName ./internal/package_name/...
  ```
- **Gate: `implementation`**
  - **PASS**: Terminal output contains `--- PASS: TestNewFeatureName` and exit code is 0.
  - **FAIL**: Test fails or times out.
