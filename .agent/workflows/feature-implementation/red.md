## 1. API Specification & Red Test (Design First)

Before writing any implementation code, prove the feature doesn't exist and define its contract.

### 1a. Write Failing Unit Test (RED)

> **Persona: SDET** — Write table-driven tests in `*_test.go` files only. Never write production code. Focus on functional and integration coverage. Cover African country/domain inputs where relevant. Tests must fail for the right reason before any implementation.

- *File*: Create or update the relevant `*_test.go` file.
- Write a test expecting the new feature to work (e.g. hitting the hypothetical endpoint).
- *Run*:
  ```bash
  // turbo
  go test -v -run TestNewFeatureName ./internal/package_name/...
  ```
- **Gate: `red-test`**
  - **PASS**: Terminal output contains `--- FAIL: TestNewFeatureName` and exit code is non-zero (specifically failing because the feature is missing).
  - **FAIL**: Test passes or fails for unrelated reasons (e.g. syntax error in test).

### 1b. Update API Specifications (Source of Truth)

> **Persona: Lead Architect** — Enforce Go best practices and clean directory structure. Coordinate TDD loop across agents. Reject non-minimal or non-modular changes. Verify spec consistency between docs and implementation.

- *Files*: Update `@[docs/api.md]` and `@[docs/openapi.yaml]` to map out the exact request/response/path for the feature. This acts as the absolute source of truth.
- Include appropriate validation rules and security requirements.
- *Lint Spec* (Optional but recommended):
  ```bash
  swagger-cli validate docs/openapi.yaml
  ```
- **Gate: `api-spec`**
  - **PASS**: Files `docs/api.md` and `docs/openapi.yaml` reflect the new feature's contract.
  - **FAIL**: Documentation is missing or inconsistent with the implementation plan.
